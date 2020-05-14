import React from 'react';
import { Loading } from 'carbon-components-react';
import Clusters from './common/Clusters';

const AppPage = ({ hasChosenAccount, tokenUpgraded, accountID }) => {
  if (!hasChosenAccount) {
    return null;
  } if (tokenUpgraded) {
    return <Clusters accountID={accountID} />;
  }
  return <Loading />;
};
export default AppPage;
