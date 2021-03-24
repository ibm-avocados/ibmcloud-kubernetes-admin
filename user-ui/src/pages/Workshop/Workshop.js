import React from 'react';
import { useParams, useLocation } from 'react-router-dom';
import Header from '../../components/Header/Header';
import {Button, TextInput} from 'carbon-components-react';

const Workshop = () => {
    const {workshop} = useParams();
    return (
        <>
            <div className="homeForm">
                <TextInput labelText="Workshop Token" id="workhop-01" type="password"/>
                <Button kind="primary">Get Cluster</Button>
            </div>
        </>
    )
}

export default Workshop;