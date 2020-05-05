import React from "react";

import {
  Form,
  TextInput,
  Button,
  Dropdown,
  TileGroup,
  RadioTile,
  Row,
  Grid,
  Column,
  FormLabel,
  Tooltip,
} from "carbon-components-react";

import styles from "./CreateForm.module.css";

import "./CreateForm.css";

const CreateForm = () => {
  const [kubernetesSelected, setKubernetesSelected] = React.useState(true);
  const [openshiftSelected, setOpenshiftSelected] = React.useState(false);

  const toggleRadio = () => {
    setKubernetesSelected(!kubernetesSelected);
    setOpenshiftSelected(!openshiftSelected);
  };

  return (
    <Form>
      <Grid>
        <FormLabel>
          <Tooltip triggerText="Cluster type and version">
            The container platform type and version for the cluster. Choose
            Kubernetes for a native Kubernetes experience on Ubuntu, or
            OpenShift to deliver Kubernetes apps on Red Hat Enterprise Linux.
          </Tooltip>
        </FormLabel>
        <Row>
          <Column md={4} lg={3}>
            <RadioTile
              onClick={toggleRadio}
              checked={kubernetesSelected}
              className={styles.radio_tile}
            >
              <div className="radio-tile-content">
                <img
                  className="radio-tile-image"
                  src="https://i.ibb.co/cDqxKBd/download.png"
                  height={100}
                  width={100}
                />
                <p>Kubernetes</p>
              </div>
              <Dropdown
                disabled={!kubernetesSelected}
                label="Select Version"
                items={["1", "2", "3"]}
              />
            </RadioTile>
          </Column>
          <Column md={4} lg={3}>
            <RadioTile
              onClick={toggleRadio}
              checked={openshiftSelected}
              className={styles.radio_tile}
            >
              <div className="radio-tile-content">
                <img
                  className="radio-tile-image"
                  src="https://i.ibb.co/0fFQCD2/openshift.png"
                  height={100}
                  width={100}
                />
                <p>OpenShift</p>
              </div>
              <Dropdown
                disabled={!openshiftSelected}
                label="Select Version"
                className={styles.dropdown}
                items={["1", "2", "3"]}
              />
            </RadioTile>
          </Column>
        </Row>
      </Grid>
    </Form>
  );
};

export default CreateForm;
