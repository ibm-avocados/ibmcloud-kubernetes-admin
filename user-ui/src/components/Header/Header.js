import React from 'react';
import {Button, Header, HeaderName} from 'carbon-components-react';
import HeaderGlobalBar from 'carbon-components-react/lib/components/UIShell/HeaderGlobalBar';
import {UserAvatar20} from '@carbon/icons-react';

import './Header.css';

function HomeHeader() {
    return (
        <>
            <Header aria-label="IBM Platform Name">
                <HeaderName href="#" prefix="IBM">
                    Developer
                </HeaderName>
                <HeaderGlobalBar>
                    <div className="userAvatar">
                        <UserAvatar20/>
                    </div>
                    <div className="userName">
                        Mofi Rahman
                    </div>
                </HeaderGlobalBar>
            </Header>
        </>
    );
}

export default HomeHeader;