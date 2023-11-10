import { NavLink } from 'react-router-dom'
import PropTypes from 'prop-types';

function SideBarItem({ title, Icon, to, href }) {
    if (href && to) {
        throw new Error('Cannot have to: and href: propos together')
    }
    return (
        <li>
            {to && <NavLink to={to} className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                {Icon && <Icon className="text-xl" />}
                <span className="ml-3">{title}</span>
            </NavLink>}
            {href && <a href={href} className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700">
                {Icon && <Icon className="text-xl" />}
                <span className="ml-3">{title}</span>
            </a>}

        </li>
    );
}

SideBarItem.propTypes = {
    title: PropTypes.string,
    to: PropTypes.string,
    Icon: PropTypes.func.isRequired,
    href: PropTypes.string
}

export default SideBarItem;


