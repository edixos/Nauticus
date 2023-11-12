import { PropTypes } from 'prop-types';

function Card({ children, ...props }) {
    return (
        <div className="bg-stone-300 shadow overflow-hidden sm:rounded-lg m-2" {...props}>
            {children}
        </div>
    );
}

Card.propTypes = {
    children: PropTypes.node.isRequired,
}


export default Card;