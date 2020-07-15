import React from 'react';
import {
  Form,
  Grid,
  Row,
  Column,
  FormLabel,
  Tooltip,
  TextInput,
  Button,
} from 'carbon-components-react';

const Spacer = ({ height }) => <div style={{ marginTop: height }} />;

const TOKEN_MESSAGE = 'github-token-saved-and-hidden';

const grab = async (url, options, retryCount = 0) => {
  const response = await fetch(url, options);
  const data = await response.json();
  if (response.status !== 200) {
    if (retryCount > 0) {
      return await grab(url, options, retryCount - 1);
    }
    throw Error(data);
  }

  return data;
};

const Settings = ({ accountID }) => {
  const [apiKeyValid, setApiKeyValid] = React.useState(false);
  const [apiKey, setApiKey] = React.useState('');
  const [org, setOrg] = React.useState('');
  const [space, setSpace] = React.useState('');
  const [region, setRegion] = React.useState('');
  const [accessGroup, setAccessGroup] = React.useState('');
  const [issueRepo, setIssueRepo] = React.useState('');
  const [grantClusterRepo, setGrantClusterRepo] = React.useState('');
  const [githubUser, setGithubUser] = React.useState('');
  const [githubToken, setGithubToken] = React.useState('');
  const [method, setMethod] = React.useState('put');
  const [data, setData] = React.useState(null);

  React.useEffect(() => {
    
    loadMetaData();
    const checkAPIKey = async () => {
      try {
        const apiKey = await fetch('/api/v1/schedule/api', {
          method: 'post',
          body: JSON.stringify({
            accountID: accountID,
          }),
        });
        if (apiKey.status === 200) {
          setApiKeyValid(true);
          setApiKey('your-api-key-will-be-pulled-from-db');
        }
      } catch (e) {
        console.log(e);
      }
    };
    checkAPIKey();
  }, [accountID]);

  
  const loadMetaData = React.useCallback(async () => {
    try {
      const metadata = await grab(`/api/v1/workshop/${accountID}/metadata`);
      if (metadata === null) {
        setMethod('post');
      } else {
        setData(metadata);
        setOrg(metadata.org);
        setSpace(metadata.space);
        setRegion(metadata.region);
        setAccessGroup(metadata.accessGroup);
        setIssueRepo(metadata.issueRepo);
        setGrantClusterRepo(metadata.grantClusterRepo);
        setGithubUser(metadata.githubUser);
        if (metadata.githubToken !== '') {
          setGithubToken(TOKEN_MESSAGE);
        }
      }
    } catch (e) {
      console.log(e);
      setMethod('post');
    }
  },[accountID]);

  const shouldUpdateBeDisabled = () => {
    if (data) {
      return (
        data.org === org &&
        data.space === space &&
        data.accessGroup === accessGroup &&
        data.githubUser === githubUser &&
        data.issueRepo === issueRepo &&
        data.grantClusterRepo === grantClusterRepo &&
        data.region === region &&
        githubToken === TOKEN_MESSAGE
      );
    }

    return (
      org === '' ||
      space === '' ||
      region === '' ||
      accessGroup === '' ||
      issueRepo === '' ||
      grantClusterRepo === '' ||
      githubUser === '' ||
      githubToken === ''
    );
  };

  const saveMetaData = async () => {
    try {
      let token = githubToken;
      if (data && token === TOKEN_MESSAGE) {
        token = data.githubToken;
      }
      const response = await grab(`/api/v1/workshop/${accountID}/metadata`, {
        method: method,
        body: JSON.stringify({
          org: org,
          space: space,
          region: region,
          accessGroup: accessGroup,
          issueRepo: issueRepo,
          grantClusterRepo: grantClusterRepo,
          githubUser: githubUser,
          githubToken: token,
        }),
      });
      loadMetaData();
    } catch (e) {
      console.log(e);
    }
  };

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
          <Button onClick={() => console.log(data)}>Save</Button>
        ) : (
          <Button kind="danger" onClick={() => console.log('delete')}>
            Delete
          </Button>
        )}
        <Spacer height="16px" />
        <Row>
          <Column sm={4} md={4} lg={3}>
            <Spacer height="16px" />
            <TextInput
              labelText="Org"
              id="account_org"
              placeholder="advowork@us.ibm.com"
              value={org}
              onChange={(e) => setOrg(e.target.value.trim())}
            />
          </Column>

          <Column sm={4} md={4} lg={3}>
            <Spacer height="16px" />
            <TextInput
              labelText="Space"
              id="account_space"
              placeholder="dev"
              value={space}
              onChange={(e) => setSpace(e.target.value.trim())}
            />
          </Column>
          <Column sm={4} md={4} lg={3}>
            <Spacer height="16px" />
            <TextInput
              labelText="Region"
              id="account_region"
              placeholder="us-south"
              value={region}
              onChange={(e) => setRegion(e.target.value.trim())}
            />
          </Column>
          <Column sm={4} md={4} lg={3}>
            <Spacer height="16px" />
            <TextInput
              labelText="Access Group"
              id="account_accessgroup"
              placeholder="lab-users"
              value={accessGroup}
              onChange={(e) => setAccessGroup(e.target.value.trim())}
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
            <Button
              disabled={shouldUpdateBeDisabled()}
              onClick={saveMetaData}
              size="default"
            >
              Update
            </Button>
          </Column>
        </Row>
      </Grid>
    </Form>
  );
};

export default Settings;
