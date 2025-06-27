// src/components/Layout.tsx
import { Outlet } from 'react-router-dom';
import Navbar from './Navbar';

const Layout = () => {
  return (
    <>
      <Navbar />
      <main className="p-4">
        <Outlet /> {/* Renders matched child routes here */}
      </main>
    </>
  );
};

export default Layout;
