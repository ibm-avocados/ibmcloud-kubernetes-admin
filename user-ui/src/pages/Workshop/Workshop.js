import React from 'react';
import { useParams, useLocation } from 'react-router-dom';

const Workshop = () => {
    const location = useLocation();
    const {workshop} = useParams();
    return (
        <div>{workshop}</div>
    )
}

export default Workshop;