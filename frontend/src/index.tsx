import { render } from "preact";
import { LocationProvider, Router, Route } from "preact-iso";
import "./style.css";

// Import all pages
import { Playground } from "./pages/playground/Playground";
import { NotFound } from "./pages/_404";
import { Login } from "./pages/auth/Login";
import { Signup } from "./pages/auth/Signup";
import { ProjectsList } from "./pages/dashboard/ProjectsList";
import { Settings } from "./pages/dashboard/Settings";
import { KanbanView } from "./pages/project/KanbanView";
import { TableView } from "./pages/project/TableView";
import { TicketDetail } from "./pages/project/TicketDetail";

export function App() {
    return (
        <LocationProvider>
            <Router>
                {/* Playground */}
                <Route path="/playground" component={Playground} />

                {/* Auth routes */}
                <Route path="/" component={Login} />
                <Route path="/signup" component={Signup} />

                {/* Dashboard routes */}
                <Route path="/dashboard/projects" component={ProjectsList} />
                <Route path="/dashboard/settings" component={Settings} />

                {/* Project routes */}
                <Route path="/dashboard/projects/:projectId/kanban" component={KanbanView} />
                <Route path="/dashboard/projects/:projectId/ticket" component={TableView} />
                <Route
                    path="/dashboard/projects/:projectId/tickets/:ticketId"
                    component={TicketDetail}
                />

                {/* 404 */}
                <Route default component={NotFound} />
            </Router>
        </LocationProvider>
    );
}

render(<App />, document.getElementById("app")!);
