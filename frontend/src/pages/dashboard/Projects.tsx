import { useEffect, useMemo, useState } from "preact/hooks";
import { useLocation } from "preact-iso";
import { AppFrame } from "../../components/layout/AppFrame";
import { APIError, api, isAbortError, Project } from "../../lib/api";
import { Badge } from "../../components/ui/Badge";
import { Button } from "../../components/ui/Button";
import { Card, CardContent } from "../../components/ui/Card";
import { Input } from "../../components/ui/Input";
import { Label } from "../../components/ui/Label";
import { Select } from "../../components/ui/Select";
import { Textarea } from "../../components/ui/Textarea";

type SortKey = "newest" | "oldest" | "name_asc" | "name_desc";

type SectionProps = {
    search: string;
    sort: SortKey;
    onSearch: (value: string) => void;
    onSort: (value: SortKey) => void;
    onCreate: () => void;
};

type CreateProjectForm = {
    name: string;
    keyPrefix: string;
    description: string;
};

function Section({ search, sort, onSearch, onSort, onCreate }: SectionProps) {
    return (
        <div class="w-full h-full flex items-center gap-2 px-2">
            <Input
                class="h-10 mt-0 flex-1 min-w-0 px-3 text-xs uppercase tracking-wider"
                placeholder="search projects"
                value={search}
                onInput={(event) => onSearch(event.currentTarget.value)}
            />

            <div class="w-36">
                <Select
                    value={sort}
                    onChange={(event) => onSort(event.currentTarget.value as SortKey)}
                    class="h-10 mt-0 px-2 text-xs uppercase tracking-wider">
                    <option value="newest">newest</option>
                    <option value="oldest">oldest</option>
                    <option value="name_asc">name a-z</option>
                    <option value="name_desc">name z-a</option>
                </Select>
            </div>

            <Button
                type="button"
                variant="primary"
                size="md"
                class="h-10 px-3 border-zinc-950 text-xs"
                onClick={onCreate}>
                create
            </Button>
        </div>
    );
}

export function Projects() {
    const { route } = useLocation();

    const [projects, setProjects] = useState<Project[]>([]);
    const [search, setSearch] = useState("");
    const [sort, setSort] = useState<SortKey>("newest");
    const [isCreateOpen, setIsCreateOpen] = useState(false);
    const [isCreating, setIsCreating] = useState(false);
    const [createError, setCreateError] = useState("");
    const [form, setForm] = useState<CreateProjectForm>({
        name: "",
        keyPrefix: "",
        description: "",
    });

    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    const visibleProjects = useMemo(() => {
        const needle = search.trim().toLowerCase();
        const filtered =
            needle ?
                projects.filter((project) => {
                    return (
                        project.name.toLowerCase().includes(needle) ||
                        project.key_prefix.toLowerCase().includes(needle) ||
                        project.description.toLowerCase().includes(needle)
                    );
                })
            :   projects;

        const sorted = [...filtered];
        sorted.sort((a, b) => {
            if (sort === "oldest") {
                return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
            }
            if (sort === "name_asc") {
                return a.name.localeCompare(b.name);
            }
            if (sort === "name_desc") {
                return b.name.localeCompare(a.name);
            }
            return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
        });

        return sorted;
    }, [projects, search, sort]);

    useEffect(() => {
        const controller = new AbortController();

        setLoading(true);
        setError("");

        const loadProjects = async () => {
            try {
                const data = await api.projects.list({ signal: controller.signal });
                if (!controller.signal.aborted) {
                    setProjects(Array.isArray(data) ? data : []);
                }
            } catch (error) {
                if (isAbortError(error)) return;

                if (error instanceof APIError) {
                    if (error.status === 401) {
                        setError("Please login to view projects.");
                    } else {
                        setError(error.message);
                    }
                } else {
                    setError("Failed to load projects.");
                }
            } finally {
                if (!controller.signal.aborted) {
                    setLoading(false);
                }
            }
        };

        loadProjects();

        return () => {
            controller.abort();
        };
    }, []);

    const createProject = async (event: Event) => {
        event.preventDefault();
        setCreateError("");

        const name = form.name.trim();
        const keyPrefix = form.keyPrefix.trim().toUpperCase();
        const description = form.description.trim();

        if (!name || !keyPrefix) {
            setCreateError("Name and key prefix are required.");
            return;
        }

        setIsCreating(true);
        try {
            const project = await api.projects.create({
                name,
                key_prefix: keyPrefix,
                description,
            });
            setProjects((prev) => [project, ...prev]);
            setForm({ name: "", keyPrefix: "", description: "" });
            setIsCreateOpen(false);
        } catch (error) {
            if (error instanceof APIError) {
                setCreateError(error.message);
            } else {
                setCreateError("Failed to create project.");
            }
        } finally {
            setIsCreating(false);
        }
    };

    return (
        <AppFrame
            section={
                <Section
                    search={search}
                    sort={sort}
                    onSearch={setSearch}
                    onSort={setSort}
                    onCreate={() => {
                        setCreateError("");
                        setIsCreateOpen((prev) => !prev);
                    }}
                />
            }>
            {isCreateOpen && (
                <form
                    class="mb-4 border border-zinc-800 bg-zinc-900 p-4 grid gap-3"
                    onSubmit={createProject}>
                    <div class="grid grid-cols-1 md:grid-cols-[1fr_180px] gap-3">
                        <div>
                            <Label htmlFor="project-name">Project name</Label>
                            <Input
                                id="project-name"
                                type="text"
                                value={form.name}
                                onInput={(event) =>
                                    setForm((prev) => ({
                                        ...prev,
                                        name: event.currentTarget.value,
                                    }))
                                }
                                placeholder="project name"
                                class="h-10"
                            />
                        </div>
                        <div>
                            <Label htmlFor="project-key">Key prefix</Label>
                            <Input
                                id="project-key"
                                type="text"
                                value={form.keyPrefix}
                                onInput={(event) =>
                                    setForm((prev) => ({
                                        ...prev,
                                        keyPrefix: event.currentTarget.value,
                                    }))
                                }
                                placeholder="key prefix"
                                maxLength={5}
                                class="h-10 uppercase"
                            />
                        </div>
                    </div>

                    <div>
                        <Label htmlFor="project-description">Description</Label>
                        <Textarea
                            id="project-description"
                            value={form.description}
                            onInput={(event) =>
                                setForm((prev) => ({
                                    ...prev,
                                    description: event.currentTarget.value,
                                }))
                            }
                            placeholder="description (optional)"
                            rows={3}
                            class="min-h-0 mt-0"
                        />
                    </div>

                    <div class="flex items-center justify-between gap-3">
                        {createError ?
                            <p class="font-mono text-xs uppercase tracking-wider text-orange-500">
                                {createError}
                            </p>
                        :   <span class="font-mono text-xs uppercase tracking-wider text-zinc-500">
                                create a new project
                            </span>
                        }
                        <div class="flex items-center gap-2">
                            <Button
                                type="button"
                                variant="outline"
                                size="sm"
                                onClick={() => setIsCreateOpen(false)}>
                                cancel
                            </Button>
                            <Button type="submit" variant="primary" size="sm" disabled={isCreating}>
                                {isCreating ? "creating..." : "create project"}
                            </Button>
                        </div>
                    </div>
                </form>
            )}
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
            : visibleProjects.length === 0 ?
                <div class="border border-zinc-800 bg-zinc-900 p-6">
                    <p class="font-mono text-xs text-zinc-400 uppercase tracking-wider">
                        {projects.length === 0 ? "No projects yet." : "No matching projects."}
                    </p>
                </div>
            :   <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                    {visibleProjects.map((project) => (
                        <Card
                            key={project.id}
                            onClick={() => route(`/dashboard/projects/${project.id}/ticket`)}
                            onKeyDown={(event) => {
                                if (event.key === "Enter" || event.key === " ") {
                                    event.preventDefault();
                                    route(`/dashboard/projects/${project.id}/ticket`);
                                }
                            }}
                            role="button"
                            tabIndex={0}
                            class="text-left p-0">
                            <CardContent class="p-4 text-zinc-400">
                                <div class="flex items-center justify-between mb-2">
                                    <Badge
                                        variant="primary"
                                        class="font-mono uppercase tracking-widest text-zinc-950 bg-lime-400/100">
                                        {project.key_prefix}
                                    </Badge>
                                    <span class="font-mono text-[10px] uppercase tracking-wider text-zinc-600">
                                        #{project.id}
                                    </span>
                                </div>
                                <h3 class="text-sm font-semibold text-zinc-100 mb-1">
                                    {project.name}
                                </h3>
                                <p class="font-mono text-[10px] text-zinc-500 leading-relaxed line-clamp-3">
                                    {project.description || "No description."}
                                </p>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            }
        </AppFrame>
    );
}
