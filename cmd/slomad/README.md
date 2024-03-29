<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# slomad

```go
import "github.com/ecshreve/slomad/cmd/slomad"
```

## Index

- [func RunDeploy\(j \*slomad.Job, confirm, force, verbose bool\) error](<#RunDeploy>)
- [func main\(\)](<#main>)
- [func newNomadClient\(\) \(\*nomadApi.Client, error\)](<#newNomadClient>)
- [func planApiJob\(nomadClient \*nomadApi.Client, job \*nomadApi.Job, diff bool\) \(bool, error\)](<#planApiJob>)
- [func submitApiJob\(nomadClient \*nomadApi.Client, job \*nomadApi.Job\) error](<#submitApiJob>)


<a name="RunDeploy"></a>
## func [RunDeploy](<https://github.com/ecshreve/slomad/blob/main/cmd/slomad/deploy.go#L15>)

```go
func RunDeploy(j *slomad.Job, confirm, force, verbose bool) error
```

RunDeploy runs a deploy for the given job.

<a name="main"></a>
## func [main](<https://github.com/ecshreve/slomad/blob/main/cmd/slomad/main.go#L12>)

```go
func main()
```



<a name="newNomadClient"></a>
## func [newNomadClient](<https://github.com/ecshreve/slomad/blob/main/cmd/slomad/deploy.go#L46>)

```go
func newNomadClient() (*nomadApi.Client, error)
```

newNomadClient returns a default nomad api client.

<a name="planApiJob"></a>
## func [planApiJob](<https://github.com/ecshreve/slomad/blob/main/cmd/slomad/deploy.go#L59>)

```go
func planApiJob(nomadClient *nomadApi.Client, job *nomadApi.Job, diff bool) (bool, error)
```

planApiJob runs a plan for the given job, and returns whether or not the job can be ignored, and any errors encountered.

<a name="submitApiJob"></a>
## func [submitApiJob](<https://github.com/ecshreve/slomad/blob/main/cmd/slomad/deploy.go#L81>)

```go
func submitApiJob(nomadClient *nomadApi.Client, job *nomadApi.Job) error
```

submitApiJob submits the given job to the nomad cluster.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
