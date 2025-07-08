import { useEffect, useState } from "react";
import axios from "axios";
/* eslint-disable  @typescript-eslint/no-explicit-any */
export default function Navbar() {
  const [logoPath, setLogoPath] = useState<string | null>(null);

  useEffect(() => {
    axios.get<{ fileName: string }>("/api/v1/logo")
      .then((res) => setLogoPath("/uploads/" + res.data.fileName))
      .catch(() => setLogoPath(null));
  }, []);

  return (
    <header className="bg-gradient-to-r from-blue-800 to-indigo-800 text-white relative">
      <div className="navbar bg-base-100 shadow-sm">
        <div className="navbar-start">
          <div className="dropdown">
            <div tabIndex={0} role="button" className="btn btn-ghost lg:hidden">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"> <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h8m-8 6h16" /> </svg>
            </div>
            <ul
              tabIndex={0}
              className="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow">
              <li><a>Attendance</a></li>
              <li>
                <a>Students</a>
                <ul className="p-2">
                  <li><a>View Students</a></li>
                  <li><a href="/students/create">Add Students</a></li>
                </ul>
              </li>
              <li><a>Item 3</a></li>
            </ul>
          </div>
          {logoPath && <img src={logoPath} alt="Logo" className="h-12 w-12" />}
          <a className="btn btn-ghost text-xl">Full Circle Martial Arts</a>
        </div>
        <div className="navbar-center hidden lg:flex">
          <ul className="menu menu-horizontal px-1">
            <li><a>Attendance</a></li>
            <li>
              <details>
                <summary>Students</summary>
                <ul className="p-2">
                  <li><a href="/students">View Students</a></li>
                  <li><a href="/students/create">Add Students</a></li>
                </ul>
              </details>
            </li>
            <li><a href="/upload-logo">Upload Logo</a></li>
          </ul>
        </div>
        <div className="navbar-end">
          <a className="btn">Button</a>
        </div>
      </div>
    </header>
  );
}
