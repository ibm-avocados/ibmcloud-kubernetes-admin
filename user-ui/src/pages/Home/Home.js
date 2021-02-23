import React from 'react';
import HomeHeader from '../../components/Header/Header';
import {Button, Form, FormGroup, RadioButton, RadioButtonGroup} from 'carbon-components-react';

import './Home.css';

const Home = () => {
    return (
        <>
            <div className="homeForm">
                <FormGroup legendText="Cluster Type">
                    <RadioButtonGroup defaultSelected="OpenShift" onChange={(e) => console.log(e)}>
                        <RadioButton labelText="OpenShift" value="OpenShift" id="radio-01"/>
                        <RadioButton labelText="Kubernetes" value="Kubernetes" id="radio-02"/>
                    </RadioButtonGroup>
                </FormGroup>
                <Button kind="primary">Get Cluster</Button>
            </div>
        </>
    )
}

export default Home;