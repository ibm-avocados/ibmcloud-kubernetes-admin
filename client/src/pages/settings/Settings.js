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
            />
          </Column>
        </Row>
        <Spacer height="16px" />
        {!apiKeyValid ? (
          <Button onClick={() => console.log("save")}>
            Save
          </Button>
        ) : (
          <Button
            kind="danger"
            onClick={() => console.log("delete")}
          >
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
            />
          </Column>

          <Column sm={4} md={4} lg={4}>
            <Spacer height="16px" />
            <TextInput labelText="Space" id="account_space" placeholder="dev" />
          </Column>
          <Column sm={4} md={4} lg={4}>
            <Spacer height="16px" />
            <TextInput
              labelText="Region"
              id="account_region"
              placeholder="us-south"
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
            />
          </Column>
          <Column sm={4} md={8} lg={6}>
            <Spacer height="16px" />
            <TextInput
              labelText="Grant Cluster URL"
              id="account_grantcluster"
              placeholder="github.ibm.com/Mofizur-Rahman/grant-cluster"
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
