import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import spacesData from '../../space.json';
import Card from "../../components/Card";
import SpaceDetailsHeader from "./components/SpaceDetailsHeader";
import SpaceDetailsBody from "./components/SpaceDetailsBody";

function SpaceDetailsPage() {
    const [spaceDetails, setSpaceDetails] = useState(null);
    const params = useParams();

    useEffect(() => {
        // Filter the space data for the space with a matching `metadata.name`
        const foundSpace = spacesData.find(space => space.metadata.name === params.spaceId);
        setSpaceDetails(foundSpace);
    }, [params.spaceId]);


    if (!spaceDetails) {
        return <div>Loading space details...</div>;
    }

    return (
        <Card>
            <SpaceDetailsHeader spaceDetails={spaceDetails} />
            <SpaceDetailsBody spaceDetails={spaceDetails} />
        </Card>
    );
}

export default SpaceDetailsPage;