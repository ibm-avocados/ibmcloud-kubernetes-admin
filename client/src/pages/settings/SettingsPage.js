import React from 'react';
import { Loading } from 'carbon-components-react';
import Settings from './Settings';

const SettingsPage = ({ hasChosenAccount, tokenUpgraded, accountID }) => {
  if (!hasChosenAccount) {
    return <h1>Please select account</h1>;
  } else if (tokenUpgraded) {
    return <Settings accountID={accountID} />;
  }
  return <Loading />;
};

export default SettingsPage;
