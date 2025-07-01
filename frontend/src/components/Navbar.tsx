


export function Navbar() {

  return (
    <header className="bg-gradient-to-r from-blue-800 to-indigo-800 text-white">
      <nav className="px-8 py-6 flex items-center relative">
        <div className="text-4xl pr-8 font-serif">Full Circle Martial Arts</div>

        <div className="relative group">
          <a
            href="#"
            className="text-lg font-semibold hover:text-blue-100 pr-8"
          >
            Students
          </a>

          <div className="absolute left-0 mt-2 hidden group-hover:block bg-white text-black shadow-lg rounded z-10">
            <a href="/students" className="block px-4 py-2 hover:bg-gray-100">
              View Students
            </a>
            <a href="/students/create" className="block px-4 py-2 hover:bg-gray-100">
              Create Student
            </a>
          </div>
        </div>

        <a className="text-lg font-semibold hover:text-blue-100 pr-8" href="/attendance">
          Attendance
        </a>

        <div className="flex-grow" />

        <div>
          <form action="/signout" method="post" className="inline pr-4">
            <button type="submit">Sign Out</button>
          </form>

          <a className="pr-4" href="/signin">Sign In</a>
          <a className="px-4 py-2 bg-blue-700 hover:bg-blue-600 rounded" href="/signup">Sign Up</a>
        </div>
      </nav>
    </header>
  );
}

export default Navbar;
