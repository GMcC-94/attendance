 
export function Navbar() {
 
  return (
     <header className="bg-gradient-to-r from-blue-800 to-indigo-800 text-white">
    <nav className="px-8 py-6 flex items-center">
      <div className="text-4xl pr-8 font-serif">Full Circle Martial Arts</div>
      <div className="">
        <a className="text-lg font-semibold hover:text-blue-100 pr-8" href="/">Home</a>
        <a className="text-lg font-semibold hover:text-blue-100 pr-8" href="/attendance">Attendance</a>
        <a className="text-lg font-semibold hover:text-blue-100 pr-8" href="/students">Students</a>
      </div>
    
      <div className="flex-grow">

      </div>

      <div>

        <form action="/signout" method="post" className="inline pr-4">
          <div className="hidden">

          </div>
          <button type="submit">Sign Out</button>
        </form>

        <a className="pr-4" href="/signin">Sign In</a>
        <a className="px-4 py-2 bg-blue-700 hover:bg-blue-600 rounded" href="/signup">Sign Up</a>

      </div>
    </nav>
  </header>

  );
};

export default Navbar;
