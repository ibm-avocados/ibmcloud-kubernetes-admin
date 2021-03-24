import React from 'react';
import { Header, HeaderName } from 'carbon-components-react';
import HeaderGlobalBar from 'carbon-components-react/lib/components/UIShell/HeaderGlobalBar';
import { UserAvatar20 } from '@carbon/icons-react';
import queryString from 'query-string';

import './Header.css';

const RootHeader = (props) => {
  const getQuery = () => {
    const { search } = location;
    const query = queryString.parse(search);
    let data = parseQuery(query.state);
    return '/auth?provider=ibm&login=true&' + data;
  }

  const parseQuery = (data) => {
    if (data === null || data === undefined) {
      return '';
    }
    return data.split('/?')[1]
  }

  return (
    <>
      <Header aria-label="IBM Platform Name">
        <HeaderName href="#" prefix="IBM">
          Developer
        </HeaderName>

        <HeaderGlobalBar>
          {props.loggedIn && props.user?
            <>
              <div className="userAvatar">
                <UserAvatar20 />
              </div>
              <div className="userName">
                {props.user.name}
              </div>
            </> : <div className="button">
              <a href={getQuery()} className="link">Sign In</a>
            </div>}
        </HeaderGlobalBar>
      </Header>
    </>
  );
};

export default RootHeader;