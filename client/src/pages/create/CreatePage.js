import React from 'react';
import { Loading } from 'carbon-components-react';
import CreateForm from './CreateForm';

const CreatePage = ({ hasChosenAccount, tokenUpgraded, accountID }) => {
  if (!hasChosenAccount) {
    return <h1>Please select account</h1>;
  } else if (tokenUpgraded) {
    return <CreateForm accountID={accountID} />;
  }
  return <Loading />;
};

export default CreatePage;
