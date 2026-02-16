import type {
	IAuthApi,
	IMembersApi,
	IProjectsApi,
	IEpicsApi,
	IBoardsApi,
	ITasksApi,
	ILabelsApi,
	IWikiApi,
	IFilesApi,
	IAdminApi
} from './types';

import { realAuthApi } from './real/auth';
import { realMembersApi } from './real/members';
import { realProjectsApi } from './real/projects';
import { realEpicsApi } from './real/epics';
import { realBoardsApi } from './real/boards';
import { realTasksApi } from './real/tasks';
import { realLabelsApi } from './real/labels';
import { realWikiApi } from './real/wiki';
import { realFilesApi } from './real/files';
import { realAdminApi } from './real/admin';

export const authApi: IAuthApi = realAuthApi;
export const membersApi: IMembersApi = realMembersApi;
export const projectsApi: IProjectsApi = realProjectsApi;
export const epicsApi: IEpicsApi = realEpicsApi;
export const boardsApi: IBoardsApi = realBoardsApi;
export const tasksApi: ITasksApi = realTasksApi;
export const labelsApi: ILabelsApi = realLabelsApi;
export const wikiApi: IWikiApi = realWikiApi;
export const filesApi: IFilesApi = realFilesApi;
export const adminApi: IAdminApi = realAdminApi;

export * from './types';
export { loadingStore, withLoading } from './loading';
