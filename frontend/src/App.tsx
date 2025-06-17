import { createBrowserRouter, Router, RouterProvider } from 'react-router-dom';
import './App.css'
import ChatSpace from './Pages/ChatSpace';


const router = createBrowserRouter([

  {
    path: "/",
    element: <ChatSpace />
  }
])

const App: React.FC = () => {

  return (

    <div>

      <RouterProvider router={router} />

    </div>

  );

};



export default App
