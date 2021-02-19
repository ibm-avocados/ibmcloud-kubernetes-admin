import React from 'react';
import {Button, Dropdown} from 'carbon-components-react';
import styles from './Navbar.module.css';
import {Logout32 as Logout, Settings32 as Settings} from '@carbon/icons-react';
// import Dropdown from "react-dropdown";
import history, { forceHistory } from '../globalHistory';
import './Dropdown.css';
import './Navbar.css';

const grab = async (url, options, retryCount = 0) => {
    const response = await fetch(url, options);
    if (response.status !== 200) {
        if (retryCount > 0) {
            console.log('failure in request. retrying again');
            return await grab(url, options, retryCount - 1);
        }
        throw Error(data);
    }
    const data = await response.json();
    return data;
};

const MenuItem = (props) => {
    return (
        <div className={props.stylesx} onClick={props.onClickHandler}>
            {props.label}
        </div>
    );
};

const MenuIcon = (props) => {
    return (
        <Button
            className="menu-icon"
            renderIcon={props.icon}
            iconDescription={props.iconDescription}
            hasIconOnly
            type="button"
            tooltipPosition="bottom"
            size="field"
            kind={props.kind}
            onClick={props.onClickHandler}
        />
    );
};

const Navbar = (props) => {
    const itemToString = (item) => {
        if (item) {
            const {name} = item.entity;
            const softlayerAccountId =
                item.entity.bluemix_subscriptions[0].softlayer_account_id || '';

            return `${name} ${softlayerAccountId}`;
        }
        return 'Unknown';
    };

    const handleCreateClick = () => {
        if (props.selectedItem) {
            history.push('/create?account=' + props.selectedItem.metadata.guid);
        } else {
            history.push('/create');
        }
    };
    const handleScheduleClick = () => {
        if (props.selectedItem) {
            history.push('/schedule?account=' + props.selectedItem.metadata.guid);
        } else {
            history.push('/schedule');
        }
    };
    const handleSettingsClick = () => {
        if (props.selectedItem) {
            history.push('/settings?account=' + props.selectedItem.metadata.guid);
        } else {
            history.push('/settings');
        }
    };

    const handleLogoutClick = async () => {
        try {
            console.log('logout clicked')
            await grab('/auth/logout', {
                method: 'post'
            })
            if (props.selectedItem) {
                forceHistory.push('/?account=' + props.selectedItem.metadata.guid);
            } else {
                forceHistory.push('/');
            }
        } catch (e) {
            console.log(e)
        }
    }

    const homeClick = () => {
        if (props.selectedItem) {
            history.push('/?account=' + props.selectedItem.metadata.guid);
        } else {
            history.push('/');
        }
    };

    return (
        <>
            <div className={styles.wrapper}>
                <div className={styles.title} onClick={homeClick}>
                    <span className={styles.bold}>IBM</span> Cloud
                </div>
                <MenuItem
                    stylesx={props.path === '/create' ? styles.activeItem : styles.item}
                    label="Create"
                    onClickHandler={handleCreateClick}
                />
                <MenuItem
                    stylesx={props.path === '/schedule' ? styles.activeItem : styles.item}
                    label="Schedule"
                    onClickHandler={handleScheduleClick}
                />
                <Dropdown
                    disabled={props.accountsLoaded}
                    className="navbar-dropdown"
                    ariaLabel="Dropdown"
                    label="Select Account"
                    items={props.items || []}
                    onChange={props.accountSelected}
                    selectedItem={props.selectedItem}
                    itemToString={itemToString}
                    id="account-dropdown"
                    light={false}
                />
                <MenuIcon
                    kind={props.path === '/settings' ? 'primary' : 'secondary'}
                    icon={Settings}
                    label="Settings"
                    onClickHandler={handleSettingsClick}
                />
                <MenuIcon
                    kind="secondary"
                    icon={Logout}
                    iconDescription="Logout"
                    label="Logout"
                    onClickHandler={handleLogoutClick}
                />

            </div>
        </>
    );
};

export default Navbar;
