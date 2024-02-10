import './App.css'
import Navbar from './Navbar'
import About from './pages/About';
import Home from './pages/Home';
import Graphs from './pages/Graphs';
import { Route, Routes } from "react-router-dom";

function App() {
  return (
    <>
      <Navbar />
      <div className="container">
        <Routes>
          <Route path='/' element={<Home />}  />
          <Route path='/graphs' element={<Graphs />}  />
          <Route path='/about' element={<About />}  />
        </Routes>
      </div>
    </>
  )
}

export default App
