import { writable, derived } from "svelte/store";
import type { ProjectResponse, CreateProjectRequest, CreateProjectsRequest, UpdateProjectRequest } from "../dtos/project-dto";
import {
    getProjects,
    getProject,
    createProject,
    createProjects,
    updateProject,
    deleteProject,
} from "../services/project-service";

function createProjectStore() {
    // ===== state =====
    const projects = writable<ProjectResponse[]>([]);
    const selectedProject = writable<ProjectResponse | null>(null);
    const searchTerm = writable("");
    const currentPage = writable(1);
    const itemsPerPage = writable(5);

    // ===== derived =====
    const filteredProjects = derived(
        [projects, searchTerm],
        ([$projects, $search]) =>
            $projects.filter((p) =>
                p.title.toLowerCase().includes($search.toLowerCase())
            )
    );

    const totalPages = derived(
        [filteredProjects, itemsPerPage],
        ([$filtered, $limit]) =>
            Math.max(1, Math.ceil($filtered.length / $limit))
    );

    const paginatedProjects = derived(
        [filteredProjects, currentPage, itemsPerPage],
        ([$filtered, $page, $limit]) => {
            const start = ($page - 1) * $limit;
            return $filtered.slice(start, start + $limit);
        }
    );

    // ===== actions =====
    function setData(data: ProjectResponse[]) {
        projects.set(data);
    }

    function setSearch(value: string) {
        searchTerm.set(value);
        currentPage.set(1);
    }

    function changePage(page: number) {
        currentPage.set(page);
    }

    function changePageSize(size: number) {
        itemsPerPage.set(size);
        currentPage.set(1);
    }

    async function fetchProjects(classroomId: string, projectRoundId: string) {
        try {
            const data = await getProjects(classroomId, projectRoundId);
            setData(data);
        } catch (err) {
            console.error("Failed to fetch projects", err);
            setData([]);
        }
    }

    async function fetchProject(classroomId: string, projectId: string) {
        try {
            const data = await getProject(classroomId, projectId);
            selectedProject.set(data);
        } catch (err) {
            console.error("Failed to fetch project", err);
            selectedProject.set(null);
        }
    }

    async function addProject(data: CreateProjectRequest) {
        try {
            const newProject = await createProject(data);
            if (!newProject || !newProject.id) {
                console.error("Dữ liệu trả về từ Service bị sai cấu trúc:", newProject);
                return;
            }
            projects.update((list) => [newProject, ...list]);
            return newProject;
        } catch (err) {
            console.error("Failed to create project", err);
            throw err;
        }
    }

    async function addProjects(data: CreateProjectsRequest) {
        try {
            const result = await createProjects(data);
            await fetchProjects(data.classroomId, data.projectRoundId); // Refresh list
            return result;
        } catch (err) {
            console.error("Failed to create projects", err);
            throw err;
        }
    }

    async function updateProjectData(data: UpdateProjectRequest) {
        try {
            const updatedProject = await updateProject(data);
            projects.update((list) =>
                list.map((p) => (p.id === data.projectId ? updatedProject : p))
            );
        } catch (err) {
            console.error("Failed to update project", err);
            throw err;
        }
    }

    async function removeProject(classroomId: string, projectId: string) {
        try {
            await deleteProject(classroomId, projectId);
            projects.update((list) => list.filter((p) => p.id !== projectId));
        } catch (err) {
            console.error("Failed to delete project", err);
            throw err;
        }
    }

    return {
        // state
        projects,
        selectedProject,
        searchTerm,
        currentPage,
        itemsPerPage,

        // derived
        filteredProjects,
        paginatedProjects,
        totalPages,

        // actions
        setData,
        setSearch,
        changePage,
        changePageSize,
        fetchProjects,
        fetchProject,
        addProject,
        addProjects,
        updateProjectData,
        removeProject,
    };
}

export const projectStore = createProjectStore();
