
import { useState } from 'react';

function Templates() {
    const [display, setDisplay] = useState(false);

    const handleDisplay = () => {
        setDisplay((prevState) => !prevState);
    }

    return (
        <div>
            {display && <h1>Toto</h1>}
            <button onClick={handleDisplay}>Display</button>
        </div>
    );
}

export default Templates;