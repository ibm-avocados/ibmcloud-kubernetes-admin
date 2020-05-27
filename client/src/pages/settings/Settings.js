import React from "react";
import {
  Form,
  Grid,
  Row,
  Column,
  FormLabel,
  Tooltip,
  TextInput,
  Button,
} from "carbon-components-react";

const Spacer = ({ height }) => <div style={{ marginTop: height }} />;

const Settings = () => {
  const [apiKeyValid, setApiKeyValid] = React.useState(false);
  const [apiKey, setApiKey] = React.useState("");
  const [org, setOrg] = React.useState("");
  const [space, setSpace] = React.useState("");
  const [region, setRegion] = React.useState("");
  const [issueRepo, setIssueRepo] = React.useState("");
  const [grantClusterRepo, setGrantClusterRepo] = React.useState("");
  const [githubUser, setGithubUser] = React.useState("");
  const [githubToken, setGithubToken] = React.useState("");

  return (
    <Form>
      <Grid>
        <FormLabel>
          <Tooltip triggerText="API Key">
            Check if your API Key is valid
          </Tooltip>
        </FormLabel>
        <Row>
          <Column sm={4} mg={8} lg={6}>
            <TextInput.PasswordInput
              labelText=""
              id="cluster_name"
              placeholder="api-key"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value.trim())}
            />
          </Column>
        </Row>
        <Spacer height="16px" />
        {!apiKeyValid ? (
          <Button onClick={() => console.log("save")}>Save</Button>
        ) : (
          <Button kind="danger" onClick={() => console.log("delete")}>
            Delete
          </Button>
        )}
        <Spacer height="16px" />
        <Row>
          <Column sm={4} md={4} lg={4}>
            <Spacer height="16px" />
            <TextInput
              labelText="Org"
              id="account_org"
              placeholder="advowork@us.ibm.com"
              value={org}
              onChange={(e) => setOrg(e.target.value.trim())}
            />
          </Column>

          <Column sm={4} md={4} lg={4}>
            <Spacer height="16px" />
            <TextInput
              labelText="Space"
              id="account_space"
              placeholder="dev"
              value={space}
              onChange={(e) => setSpace(e.target.value.trim())}
            />
          </Column>
          <Column sm={4} md={4} lg={4}>
            <Spacer height="16px" />
            <TextInput
              labelText="Region"
              id="account_region"
              placeholder="us-south"
              value={region}
              onChange={(e) => setRegion(e.target.value.trim())}
            />
          </Column>
        </Row>

        <Row>
          <Column sm={4} md={8} lg={6}>
            <Spacer height="16px" />
            <TextInput
              labelText="Github Issue Repo"
              id="account_gitrepo"
              placeholder="github.ibm.com/jja/cloud-workshop-requests"
              value={issueRepo}
              onChange={(e) => setIssueRepo(e.target.value.trim())}
            />
          </Column>
          <Column sm={4} md={8} lg={6}>
            <Spacer height="16px" />
            <TextInput
              labelText="Grant Cluster URL"
              id="account_grantcluster"
              placeholder="github.ibm.com/Mofizur-Rahman/grant-cluster"
              value={grantClusterRepo}
              onChange={(e) => setGrantClusterRepo(e.target.value.trim())}
            />
          </Column>
        </Row>

        <Row>
          <Column sm={4} md={8} lg={6}>
            <Spacer height="16px" />
            <TextInput
              labelText="Github Username"
              id="account_gituser"
              placeholder="Mofizur-Rahman"
              value={githubUser}
              onChange={(e) => setGithubUser(e.target.value.trim())}
            />
          </Column>
          <Column sm={4} md={8} lg={6}>
            <Spacer height="16px" />
            <TextInput.PasswordInput
              labelText="Github Token"
              id="account_github_token"
              placeholder="your-token-here"
              value={githubToken}
              onChange={(e) => setGithubToken(e.target.value.trim())}
            />
          </Column>
        </Row>
        <Spacer height="16px" />
        <Row>
          <Column>
            <Button onClick={() => console.log("save")} size="default">
              Update
            </Button>
          </Column>
        </Row>
      </Grid>
    </Form>
  );
};

export default Settings;
