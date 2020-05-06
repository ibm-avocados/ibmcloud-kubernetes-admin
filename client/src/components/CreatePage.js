import React from "react";
import CreateForm from "./common/CreateForm";

const CreatePage = ({ accountID, hasChosenAccount }) => {
  if (hasChosenAccount) {
    return <CreateForm accountID={accountID} />;
  } else {
    return <h1>Please select account</h1>;
  }
};

export default CreatePage;
