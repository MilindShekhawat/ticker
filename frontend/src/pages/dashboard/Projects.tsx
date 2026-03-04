import { useEffect, useState } from "preact/hooks";
import { useLocation } from "preact-iso";
import { AppFrame } from "../../components/layout/AppFrame";
import { APIError, api, isAbortError, Project } from "../../lib/api";
import { Button } from "../../components/ui/Button";

function Section() {
    return (
        <div class="px-3">
            <p class="text-xs uppercase tracking-wider text-zinc-300">Projects</p>
        </div>
    );
}

export function Projects() {
    const { route } = useLocation();

    const [projects, setProjects] = useState<Project[]>([]);

    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        const controller = new AbortController();

        setLoading(true);
        setError("");

        (async () => {
            try {
                const data = await api.projects.list({ signal: controller.signal });
                setProjects(data);
            } catch (error) {
                if (isAbortError(error)) return;

                if (error instanceof APIError) {
                    // TODO: Send proper message from backend
                    if (error.status === 401) setError("Please login to view projects.");
                    else setError(error.message);
                } else {
                    setError("Failed to load projects.");
                }
            } finally {
                console.log("WHAT IS LOADING");
              if (!controller.signal.aborted) {
                console.log("WHAT IS LOADING Fase");
                setLoading(false);
              }
            }
        })();

        return () => {
            controller.abort();
        };
    }, []);

    return (
        <AppFrame section={<Section />}>
            {loading ?
                <div class="border border-zinc-800 bg-zinc-900 p-6">
                    <p class="font-mono text-xs text-zinc-400 uppercase tracking-wider">
                        Loading projects...
                    </p>
                </div>
            : error ?
                <div class="border border-zinc-800 bg-zinc-900 p-6 flex items-center justify-between gap-4">
                    <p class="font-mono text-xs text-orange-500 uppercase tracking-wider">
                        {error}
                    </p>
                    <Button variant="outline" size="sm" onClick={() => route("/")}>
                        login
                    </Button>
                </div>
            : projects.length === 0 ?
                <div class="border border-zinc-800 bg-zinc-900 p-6">
                    <p class="font-mono text-xs text-zinc-400 uppercase tracking-wider">
                        No projects yet.
                    </p>
                </div>
            :   <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                    {projects.map((project) => (
                        <button
                            key={project.id}
                            type="button"
                            onClick={() => route(`/dashboard/projects/${project.id}/ticket`)}
                            class="text-left border border-zinc-800 bg-zinc-900 hover:border-lime-400 transition p-4">
                            <div class="flex items-center justify-between mb-2">
                                <span class="font-mono text-xs font-bold tracking-widest text-zinc-950 bg-lime-400 px-2 py-0.5 uppercase">
                                    {project.key_prefix}
                                </span>
                                <span class="font-mono text-[10px] uppercase tracking-wider text-zinc-600">
                                    #{project.id}
                                </span>
                            </div>
                            <h3 class="text-sm font-semibold text-zinc-100 mb-1">{project.name}</h3>
                            <p class="font-mono text-[10px] text-zinc-500 leading-relaxed line-clamp-3">
                                {project.description || "No description."}
                            </p>
                        </button>
                    ))}
                </div>
            }
        </AppFrame>
    );
}
