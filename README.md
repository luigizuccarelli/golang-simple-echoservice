# A simple golang microservice for OC4 Tekton POC

Uses a simple post endpoint to echo payload message data.

## Workflow

I make use of TDD (Trunk driven development) as suggested by Gene Kim et. al. in the book **Accelerate**. It allows for fast time to market CICD builds, deploy and delivery (highly recommended read).

The process is as follows (assume repos's have been created and initial code has been committed)
- I use 2 repo's to accomplish this
- Bitbucket stash (source of truth)
- Gitea hosted on Openshift (staging/work repo)

**NB** I have left all the pipeline artifacts (json files, shell scripts etc) in the repo so that the example can be referenced easily for porting to tekton

#### Step 1 
- Checkout master from Bitbucket (stash) (SOT)
- Ensure that both gitea and stash (remote settings exist) - I use 'origin' for gitea and 'stash' for bitbucket

#### Step 2 - Create a PR on Gitea 
- Make changes to the code
- Create a PR on Gitea ```git push origin <dev-branch>```
- Ask for a peer revue and verification 
  - Reviewer checks out the PR
  - Execute clean  ```make clean```
  - Execute tests  ```make test```
  - Execute cover  ```make cover```
  - Execute build  ```make build``` (compiles the code)

- If all goes well the reviewer gives the thumbs up to merge

#### Step 3 - Merge PR
- This triggers the GOCD pipeline
  - Checks out master from Gitea
  - Executes ```make clean```
  - Executes ```make test```
  - Executes ```make cover```
  - Executes ```make build``` 
  - Executes gocd shell script - sonarcube scanner
  - Executes gocd shell script - kaniko build -> build and push image to JFrog

#### Step 4 - Deploy on Openshift (manual process)
- Login to openshift
- use ```oc rollout --latest```
- monitor and check that the deploy is working
- ```oc get pods```

#### Step 5 - Check endpoints
- Use curl or other (FE) to check API calls work

#### Step 6 - Merge to Bitbucket (manual process)
- ```git push stash master```

## Testing locally
- Execute ```make clean && make build```

- Execute ```./run.sh``` (this will setup the correct envars and execute the microservice binary)

- ```curl -d '{"id":"123456","message":"Hello World !!!!"}' http://127.0.0.1:9000/api/v1/echo```

- You should get the following repsonse :

- ```{"code": 200,"status": "OK",	"message": "Hello World !!!!"}```

- Endpoint for **OpenAPI 3.0** (swagger docs) - ```http://127.0.0.1:9000/api/v1/api-docs/``` (trailing slash is important)


## Links
Gitea     - https://gitea-cicd.apps.aws2-dev.ocp.14west.io/

GOCD      - https://gocdserver-cicd.apps.aws2-dev.ocp.14west.io/

SonarQube - http://ec2-3-248-27-131.eu-west-1.compute.amazonaws.com/projects

Detailed CICD flow and architecture - https://wiki.14west.us/pages/viewpage.action?pageId=93263322 (section **CICD Pipeline Architecture**)
 
I will send the credentials if needed on request
