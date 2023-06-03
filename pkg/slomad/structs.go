package slomad

import (
	"fmt"
	"os"
	"sort"

	"github.com/google/uuid"
	nomadApi "github.com/hashicorp/nomad/api"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

// CreateApiJob returns a validated nomad api job. The `force` argument
// controls whether or not we want to explicitly force a new version.
func (j *Job) CreateApiJob(force bool) (*nomadApi.Job, error) {
	pl := struct{ Labels []string }{}
	config := make(map[string]interface{})
	if j.Image != nil {
		config["image"] = j.Image
	}
	if j.Ports != nil {
		config["ports"] = pl.Labels
	}
	if j.Args != nil {
		config["args"] = j.Args
	}
	if j.Volumes != nil {
		vols := []string{}
		for src, dst := range j.Volumes {
			vols = append(vols, fmt.Sprintf("%s:%s", src, dst))
		}
		sort.Strings(vols)
		config["volumes"] = vols
	}
	if j.Caps != nil {
		config["cap_add"] = j.Caps
	}

	if j.Name == "glances" {
		config["network_mode"] = "host"
	}

	var pluginConfig = &nomadStructs.TaskCSIPluginConfig{}
	csiVols := make(map[string]*nomadStructs.VolumeRequest)
	if j.Storage != nil {
		if *j.Storage == "controller" || *j.Storage == "node" {
			config["privileged"] = true
			config["network_mode"] = "host"
			pluginConfig = &nomadStructs.TaskCSIPluginConfig{
				ID:       "nfs",
				Type:     nomadStructs.CSIPluginType(*j.Storage),
				MountDir: "/csi",
			}
		} else {
			pluginConfig = nil
			volName := fmt.Sprintf("%s-vol", *j.Storage)
			csiVols[volName] = &nomadStructs.VolumeRequest{
				Name:   volName,
				Source: volName,

				Type:           "csi",
				ReadOnly:       false,
				AccessMode:     "single-node-writer",
				AttachmentMode: "file-system",
			}
		}
	} else {
		pluginConfig = nil
	}

	resources := &nomadStructs.Resources{
		CPU:      j.Size["cpu"],
		MemoryMB: j.Size["mem"],
	}

	constraint := &nomadStructs.Constraint{
		LTarget: "${attr.unique.hostname}",
		RTarget: j.Constraint,
		Operand: "regexp",
	}

	services := []*nomadStructs.Service{}
	for _, pl := range pl.Labels {
		services = append(services, GetService(j.Name, pl))
	}

	templates := []*nomadStructs.Template{}
	if j.Templates != nil {
		for tmplname, tmpl := range j.Templates {
			templates = append(templates, &nomadStructs.Template{
				EmbeddedTmpl: tmpl,
				DestPath:     fmt.Sprintf("local/config/%s", tmplname),
				ChangeMode:   "signal",
				ChangeSignal: "SIGHUP",
			})
		}
	}

	user := ""
	if j.User != nil {
		user = *j.User
	}

	volumeMounts := []*nomadStructs.VolumeMount{}
	if j.Mounts != nil {
		for src, dst := range j.Mounts {
			volumeMounts = append(volumeMounts, &nomadStructs.VolumeMount{
				Volume:      src,
				Destination: dst,
				ReadOnly:    false,
			})
		}
	}
	task := &nomadStructs.Task{
		Name:            j.Name,
		Driver:          j.Driver,
		Config:          config,
		User:            user,
		Resources:       resources,
		Env:             j.Env,
		Templates:       templates,
		LogConfig:       nomadStructs.DefaultLogConfig(),
		Services:        services,
		CSIPluginConfig: pluginConfig,
		VolumeMounts:    volumeMounts,
	}

	reschedulePolicy := &nomadStructs.DefaultServiceJobReschedulePolicy
	if j.JobType == "system" {
		reschedulePolicy = nil
	}

	group := &nomadStructs.TaskGroup{
		Name:             j.Name,
		Count:            j.Count,
		Tasks:            []*nomadStructs.Task{task},
		RestartPolicy:    nomadStructs.NewRestartPolicy("service"),
		ReschedulePolicy: reschedulePolicy,
		Volumes:          csiVols,
		EphemeralDisk: &nomadStructs.EphemeralDisk{
			SizeMB: 500,
		},
		Networks: []*nomadStructs.NetworkResource{
			{
				ReservedPorts: nil,
				DynamicPorts:  nil,
			},
		},
	}

	jobType := "service"
	if j.JobType != "" {
		jobType = j.JobType
	}

	job := &nomadStructs.Job{
		ID:          j.Name,
		Name:        j.Name,
		Region:      "global",
		Priority:    j.Priority,
		Datacenters: []string{"dcs"},
		Type:        jobType,
		TaskGroups:  []*nomadStructs.TaskGroup{group},
		Namespace:   "default",
		Constraints: []*nomadStructs.Constraint{constraint},
	}

	// Writing a new uuid to this field ensures Nomad will create a new
	// version of the job.
	if force {
		job.Meta = map[string]string{
			"run_uuid": uuid.NewString(),
		}
	}

	if err := job.Validate(); err != nil {
		log.Errorf("Nomad job validation failed. Error: %s\n", err)
		return nil, err
	}

	apiJob, err := convertJob(job)
	if err != nil {
		log.Errorf("Failed to convert nomad job in api call. Error: %s\n", err)
		return nil, err
	}

	return apiJob, nil
}

// planApiJob creates a nomad api client, and runs a plan
// for the given job, printing the output diff.
func planApiJob(job *nomadApi.Job) error {
	nomadConfig := nomadApi.DefaultConfig()
	nomadConfig.Address = os.Getenv("NOMAD_TARGET")
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return oops.Wrapf(err, "unable to create nomad api client")
	}

	planResp, _, nomadErr := nomadClient.Jobs().Plan(job, true, nil)
	if nomadErr != nil {
		log.Errorf("Error planning job: %s", nomadErr)
		return fmt.Errorf(fmt.Sprintf("Error planning job: %s", nomadErr))
	}

	log.Infof("Sucessfully planned nomad job %s - %+v\n", *job.Name, planResp.Annotations.DesiredTGUpdates[*job.Name])
	return nil
}

// submitApiJob creates a nomad api client, and submits the job to nomad.
func submitApiJob(job *nomadApi.Job) error {
	nomadConfig := nomadApi.DefaultConfig()
	nomadConfig.Address = os.Getenv("NOMAD_TARGET")
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return oops.Wrapf(err, "unable to create nomad api client")
	}

	_, _, nomadErr := nomadClient.Jobs().Register(job, nil)
	if nomadErr != nil {
		return oops.Wrapf(err, "error submitting job: %+v", job)
	}

	log.Infof("Sucessfully submitted nomad job %s\n", *job.Name)
	return nil
}

// Plan creates a new API job and runs a plan on it.
func (j *Job) Plan(force bool) error {
	apiJob, err := j.CreateApiJob(force)
	if err != nil {
		return oops.Wrapf(err, "error crating api job for job: %+v", j)
	}

	if err = planApiJob(apiJob); err != nil {
		return oops.Wrapf(err, "error planning api job")
	}

	return nil
}

// Plan creates a new API job and runs a plan on it.
func (j *App) PlanApp(force bool) error {
	_, aj, err := j.ToNomadJob(force)
	if err != nil {
		return oops.Wrapf(err, "error creating api job for job: %+v", j)
	}

	if err = planApiJob(aj); err != nil {
		return oops.Wrapf(err, "error planning api job")
	}

	return nil
}

// Deploy creates a new ApiJob and submits it to the nomad API.
func (j *Job) Deploy(force bool) error {
	apiJob, err := j.CreateApiJob(force)
	if err != nil {
		return oops.Wrapf(err, "error creating api job for job: %+v", j)
	}

	if err = submitApiJob(apiJob); err != nil {
		return oops.Wrapf(err, "error submitting api job")
	}

	return nil
}

func (j *App) DeployApp(force bool) error {
	_, aj, err := j.ToNomadJob(force)
	if err != nil {
		return oops.Wrapf(err, "error creating api job for job: %+v", j)
	}

	if err = submitApiJob(aj); err != nil {
		return oops.Wrapf(err, "error submitting api job")
	}

	return nil
}
