import type {
	User,
	Community,
	Project,
	Epic,
	Board,
	Task,
	Label,
	WikiPage,
	Attachment,
	DashboardStats,
	CommunitySettings
} from '$lib/types/api';

export interface ApiResponse<T> {
	data: T | null;
	error: string | null;
	status: number;
}

export interface PaginatedResponse<T> {
	items: T[];
	total: number;
	page: number;
	pageSize: number;
	totalPages: number;
}

export interface ApiClient {
	get<T>(url: string): Promise<ApiResponse<T>>;
	post<T>(url: string, body?: unknown): Promise<ApiResponse<T>>;
	put<T>(url: string, body?: unknown): Promise<ApiResponse<T>>;
	patch<T>(url: string, body?: unknown): Promise<ApiResponse<T>>;
	delete<T>(url: string): Promise<ApiResponse<T>>;
}

export class ApiError extends Error {
	constructor(
		public status: number,
		public code: string,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

export function isApiError(error: unknown): error is ApiError {
	return error instanceof ApiError;
}

export function isUnauthorized(error: unknown): boolean {
	return isApiError(error) && error.status === 401;
}

export function isNotFound(error: unknown): boolean {
	return isApiError(error) && error.status === 404;
}

export function isConflict(error: unknown): boolean {
	return isApiError(error) && error.status === 409;
}

export function createResponse<T>(data: T): ApiResponse<T> {
	return { data, error: null, status: 200 };
}

export function createErrorResponse<T>(error: string, status: number): ApiResponse<T> {
	return { data: null, error, status };
}

export interface RegisterInput {
	communityName: string;
	slug: string;
	adminEmail: string;
	adminPassword: string;
}

export interface InviteMemberInput {
	email: string;
	role: 'admin' | 'member';
}

export interface UpdateProfileInput {
	firstName?: string;
	lastName?: string;
}

export interface CreateProjectInput {
	name: string;
	description: string;
}

export interface UpdateProjectInput {
	name?: string;
	description?: string;
	status?: 'active' | 'archived';
}

export interface CreateEpicInput {
	name: string;
	description: string;
	color: string;
	projectId: string;
}

export interface UpdateEpicInput {
	name?: string;
	description?: string;
	color?: string;
}

export interface CreateBoardInput {
	name: string;
	description?: string;
	projectId: string;
}

export interface CreateTaskInput {
	title: string;
	description?: string;
	boardId: string;
	epicId?: string;
	assigneeId?: string;
	labels?: string[];
}

export interface UpdateTaskInput {
	title?: string;
	description?: string;
	status?: 'backlog' | 'todo' | 'inprogress' | 'done';
	assigneeId?: string;
	epicId?: string;
	labels?: string[];
}

export interface MoveTaskInput {
	status: 'backlog' | 'todo' | 'inprogress' | 'done';
	position: number;
}

export interface CreateLabelInput {
	name: string;
	color: string;
	communityId: string;
}

export interface CreateWikiPageInput {
	title: string;
	slug: string;
	content: Record<string, unknown>;
	projectId: string;
	parentId?: string;
}

export interface UpdateWikiPageInput {
	title?: string;
	slug?: string;
	content?: Record<string, unknown>;
	parentId?: string;
}

export interface UploadFileOptions {
	taskId?: string;
	wikiPageId?: string;
	uploadedBy: string;
}

export interface ListAttachmentsOptions {
	taskId?: string;
	wikiPageId?: string;
}

export interface UpdateCommunitySettingsInput {
	name?: string;
	emailConfig?: CommunitySettings['emailConfig'];
	storageConfig?: CommunitySettings['storageConfig'];
	metricsEnabled?: boolean;
}

export interface IAuthApi {
	login(email: string, password: string): Promise<ApiResponse<{ user: User; token: string }>>;
	logout(): Promise<ApiResponse<void>>;
	getCurrentUser(): Promise<ApiResponse<User>>;
	register(input: RegisterInput): Promise<ApiResponse<{ user: User; community: Community; token: string }>>;
}

export interface IMembersApi {
	listMembers(): Promise<ApiResponse<PaginatedResponse<User>>>;
	inviteMember(input: InviteMemberInput): Promise<ApiResponse<{ invitationId: string }>>;
	updateProfile(userId: string, input: UpdateProfileInput): Promise<ApiResponse<User>>;
}

export interface IProjectsApi {
	listProjects(): Promise<ApiResponse<PaginatedResponse<Project>>>;
	getProject(id: string): Promise<ApiResponse<Project>>;
	createProject(input: CreateProjectInput): Promise<ApiResponse<Project>>;
	updateProject(id: string, input: UpdateProjectInput): Promise<ApiResponse<Project>>;
	deleteProject(id: string): Promise<ApiResponse<void>>;
}

export interface IEpicsApi {
	listEpics(projectId: string): Promise<ApiResponse<Epic[]>>;
	createEpic(input: CreateEpicInput): Promise<ApiResponse<Epic>>;
	updateEpic(id: string, input: UpdateEpicInput): Promise<ApiResponse<Epic>>;
	deleteEpic(id: string): Promise<ApiResponse<void>>;
}

export interface IBoardsApi {
	listBoards(projectId: string): Promise<ApiResponse<Board[]>>;
	getBoard(id: string): Promise<ApiResponse<Board>>;
	createBoard(input: CreateBoardInput): Promise<ApiResponse<Board>>;
	deleteBoard(id: string): Promise<ApiResponse<void>>;
}

export interface ITasksApi {
	listTasks(boardId: string): Promise<ApiResponse<Task[]>>;
	createTask(input: CreateTaskInput): Promise<ApiResponse<Task>>;
	updateTask(id: string, input: UpdateTaskInput): Promise<ApiResponse<Task>>;
	moveTask(id: string, input: MoveTaskInput): Promise<ApiResponse<Task>>;
	deleteTask(id: string): Promise<ApiResponse<void>>;
}

export interface ILabelsApi {
	listLabels(communityId: string): Promise<ApiResponse<Label[]>>;
	createLabel(input: CreateLabelInput): Promise<ApiResponse<Label>>;
	assignLabelToTask(taskId: string, labelId: string): Promise<ApiResponse<Task>>;
}

export interface IWikiApi {
	listWikiPages(projectId: string): Promise<ApiResponse<WikiPage[]>>;
	getWikiPage(id: string): Promise<ApiResponse<WikiPage>>;
	createWikiPage(input: CreateWikiPageInput): Promise<ApiResponse<WikiPage>>;
	updateWikiPage(id: string, input: UpdateWikiPageInput): Promise<ApiResponse<WikiPage>>;
	deleteWikiPage(id: string): Promise<ApiResponse<void>>;
}

export interface IFilesApi {
	uploadFile(file: File, options: UploadFileOptions): Promise<ApiResponse<Attachment>>;
	downloadFile(id: string): Promise<ApiResponse<Blob>>;
	listAttachments(options: ListAttachmentsOptions): Promise<ApiResponse<Attachment[]>>;
}

export interface IAdminApi {
	getCommunitySettings(): Promise<ApiResponse<CommunitySettings>>;
	updateCommunitySettings(input: UpdateCommunitySettingsInput): Promise<ApiResponse<CommunitySettings>>;
	getDashboardStats(): Promise<ApiResponse<DashboardStats>>;
}
