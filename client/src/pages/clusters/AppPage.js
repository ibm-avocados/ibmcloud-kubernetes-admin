import React from 'react';
import { Loading } from 'carbon-components-react';
import Clusters from './Clusters';

const AppPage = ({ hasChosenAccount, tokenUpgraded, accountID }) => {
  if (!hasChosenAccount) {
    return null;
  } else if (tokenUpgraded) {
    return <Clusters accountID={accountID} />;
  }
  return <Loading />;
};
export default AppPage;
