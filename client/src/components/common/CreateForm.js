import React from "react";

import {
  Form,
  TextInput,
  Button,
  Dropdown,
  TileGroup,
  RadioTile
} from "carbon-components-react";

import "./CreateForm.css";

const CreateForm = () => {
  return (
    <Form>
      <TextInput id="1" labelText="1" />
      <TextInput id="2" labelText="2" />
      <TextInput id="3" labelText="3" />
      <TextInput id="4" labelText="4" />
      <TextInput id="5" labelText="5" />
      <TextInput id="6" labelText="6" />
      <TextInput id="7" labelText="7" />
    </Form>
  );
};

export default CreateForm;
