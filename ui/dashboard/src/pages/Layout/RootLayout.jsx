// RootLayout.js
import { Outlet } from 'react-router-dom';
import SideBar from '../../components/SideBar';
import Header from '../../components/Header';
import Logo from '../../assets/nauticus-logo.png'
const RootLayout = () => {
    return (
        <div className="flex flex-col h-screen bg-gray-100">
            <Header title="Nauticus" logoSrc={Logo} logoAlt="Nauticus Dashboard Logo" />
            <div className="flex flex-1 overflow-hidden">
                <SideBar />
                <main className="flex-1 overflow-y-auto bg-gray-600">
                    <Outlet /> {/* This is where the child routes will render */}
                </main>
            </div>
        </div>
    );
};

export default RootLayout;