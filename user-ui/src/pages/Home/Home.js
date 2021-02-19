import React from 'react';
import HomeHeader from '../../components/Header';
import {Button, Form, FormGroup, RadioButton, RadioButtonGroup} from 'carbon-components-react';

const Home = () => {
    return (
        <>
            <HomeHeader/>
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