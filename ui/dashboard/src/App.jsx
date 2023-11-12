import './App.css'

import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import RootLayout from './pages/Layout/RootLayout';
import SpacesPage from './pages/Spaces/SpacesPage';
import Templates from './pages/Templates/Templates';
import SpaceDetailsPage from './pages/Spaces/SpaceDetailsPage';


const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    children: [
      {
        index: true, element: <SpacesPage />
      },
      {
        path: 'spaces', children: [
          { index: true, element: <SpacesPage /> },
          { path: ':spaceId', element: <SpaceDetailsPage /> }
        ]
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