import React from "react";

const Navbar: React.FC = () => {
  return (
    

<div className="w-full bg-blue-700 p-4 text-white flex items-center justify-start space-x-6">
  <div className="font-bold text-xl">My App</div>
  <nav className="flex space-x-4">
    <a href="#" className="hover:underline">Home</a>
    <a href="#" className="hover:underline">Students</a>
    <a href="#" className="hover:underline">Attendance</a>
  </nav>
</div>


  );

};

export default Navbar;
