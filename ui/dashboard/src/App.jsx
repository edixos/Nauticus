import './App.css'

import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import RootLayout from './pages/Layout/RootLayout';
import Spaces from './pages/Spaces/Spaces';
import Templates from './pages/Templates/Templates';


const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    children: [
      {
        index: true, element: <Spaces />
      },
      {
        path: 'spaces', element: <Spaces />
      },
      {
        path: 'templates', element: <Templates />
      },
    ]
  }
])

function App() {
  return <RouterProvider router={router} />;
}

export default App;