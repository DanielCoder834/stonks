import './index.css'
import { Link, useMatch, useResolvedPath } from "react-router-dom";

export default function Navbar() {
    return <>
    <nav className="nav">
        <Link to="/" className="site-title">Stonks</Link>
        <ul id='nav-list'>
            <CustomLink to="/graphs" id="pricingLink">Graphs</CustomLink>
            <CustomLink to="/about" id="aboutLink">About</CustomLink>
        </ul>
    </nav>
    </>
}

function CustomLink({ to , children, id, ...props }: {to:any, children: any, id: any}) {
    const resolvedPath = useResolvedPath(to)
    const isActive = useMatch({ path: resolvedPath.pathname, end:true })
    return (
        <li className={isActive ? "active" : ""} id={id}>
            <Link to={to} {... props}>{children}</Link>
        </li>
    )
}
