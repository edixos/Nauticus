import { IoDocumentText } from 'react-icons/io5'
import { HiOutlineCube, HiCubeTransparent } from 'react-icons/hi2'
import SideBarItem from './SideBarItem';
// Import icons from 'react-icons' as needed

function SideBar() {
    return (
        <aside className="w-64 h-screen bg-gray-50 dark:bg-gray-800 px-4 py-6">
            <nav>
                <ul>
                    <SideBarItem title="Spaces" Icon={HiOutlineCube} to='/spaces' />
                    <SideBarItem title="Templates" Icon={HiCubeTransparent} to='/templates' />
                    <SideBarItem title="Documentation" Icon={IoDocumentText} href='https://nauticus.edixos.com/0.1.3/' />
                    {/* Add more SideBarItem as needed */}
                </ul>
            </nav>
        </aside>
    );
}

export default SideBar;
