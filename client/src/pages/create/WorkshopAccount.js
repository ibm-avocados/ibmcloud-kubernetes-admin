import React from 'react';
import {
  TextInput,
  Checkbox,
  InlineLoading,
  InlineNotification,
  Button,
} from 'carbon-components-react';

import history from '../../globalHistory';

const Spacer = ({ height }) => <div style={{ marginTop: height }} />;

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

const WorkshopAccount = ({ accountID, githubIssue, setGithubIssue, isWorkshop, setIsWorkshop }) => {
  return (
    <>
      <Checkbox
        id="workshop_checkbox"
        labelText="Workshop"
        onChange={(e) => setIsWorkshop(e)}
      />
      {isWorkshop ? (
        <WorkshopView
          githubIssue={githubIssue}
          setGithubIssue={setGithubIssue}
          accountID={accountID}
        />
      ) : (
        <></>
      )}
    </>
  );
};

const WorkshopView = ({ accountID, setGithubIssue, githubIssue }) => {
  const [loading, setLoading] = React.useState(false);
  const [metadataAvailable, setMetadataAvailable] = React.useState(false);
  const [apiKeyValid, setApiKeyValid] = React.useState(false);

  const onSettingsButtonClicked = () => {
    history.push('/settings');
  };
  React.useEffect(() => {
    const checkMetadata = async () => {
      setLoading(true);
      try {
        const metadata = await grab(`/api/v1/workshop/${accountID}/metadata`);
        if (metadata !== null) {
          setMetadataAvailable(true);
        }
        const apiKey = await fetch('/api/v1/schedule/api', {
          method: 'post',
          body: JSON.stringify({
            accountID: accountID,
          }),
        });
        if (apiKey.status === 200) {
          setApiKeyValid(true);
        }
      } catch (e) {
        console.log(e);
      }
      setLoading(false);
    };
    checkMetadata();
  }, [accountID]);

  if (loading) {
    return <InlineLoading description="checking information" />;
  }

  if (!apiKeyValid || !metadataAvailable) {
    return (
      <>
        <Button onClick={onSettingsButtonClicked} kind="primary" size="small">
          Settings
        </Button>
        <Spacer height="32px" />
      </>
    );
  }
  return (
    <TextInput
      id="workshop_issue"
      value={githubIssue}
      invalid={isNaN(githubIssue) || githubIssue === ''}
      invalidText="Must be a number"
      onChange={(e) => setGithubIssue(e.target.value.trim())}
      labelText="Github Issue Number"
    />
  );
};

export default WorkshopAccount;
