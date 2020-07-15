import React from 'react';
import { Loading } from 'carbon-components-react';
import Clusters from './Clusters';

const AppPage = (props) => {
  const { hasChosenAccount, tokenUpgraded, accountID, query } = props;
  if (!hasChosenAccount) {
    return null;
  } else if (tokenUpgraded) {
    return <Clusters query={query} accountID={accountID} />;
  }
  return <Loading />;
};
export default AppPage;
