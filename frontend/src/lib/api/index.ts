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

import { mockAuthApi } from './mock/auth';
import { mockMembersApi } from './mock/members';
import { mockProjectsApi } from './mock/projects';
import { mockEpicsApi } from './mock/epics';
import { mockBoardsApi } from './mock/boards';
import { mockTasksApi } from './mock/tasks';
import { mockLabelsApi } from './mock/labels';
import { mockWikiApi } from './mock/wiki';
import { mockFilesApi } from './mock/files';
import { mockAdminApi } from './mock/admin';

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

const USE_MOCK = import.meta.env.PUBLIC_USE_MOCK_API !== 'false';

let authApiImpl: IAuthApi;
let membersApiImpl: IMembersApi;
let projectsApiImpl: IProjectsApi;
let epicsApiImpl: IEpicsApi;
let boardsApiImpl: IBoardsApi;
let tasksApiImpl: ITasksApi;
let labelsApiImpl: ILabelsApi;
let wikiApiImpl: IWikiApi;
let filesApiImpl: IFilesApi;
let adminApiImpl: IAdminApi;

if (USE_MOCK) {
	authApiImpl = mockAuthApi;
	membersApiImpl = mockMembersApi;
	projectsApiImpl = mockProjectsApi;
	epicsApiImpl = mockEpicsApi;
	boardsApiImpl = mockBoardsApi;
	tasksApiImpl = mockTasksApi;
	labelsApiImpl = mockLabelsApi;
	wikiApiImpl = mockWikiApi;
	filesApiImpl = mockFilesApi;
	adminApiImpl = mockAdminApi;
} else {
	authApiImpl = realAuthApi;
	membersApiImpl = realMembersApi;
	projectsApiImpl = realProjectsApi;
	epicsApiImpl = realEpicsApi;
	boardsApiImpl = realBoardsApi;
	tasksApiImpl = realTasksApi;
	labelsApiImpl = realLabelsApi;
	wikiApiImpl = realWikiApi;
	filesApiImpl = realFilesApi;
	adminApiImpl = realAdminApi;
}

export const authApi = authApiImpl;
export const membersApi = membersApiImpl;
export const projectsApi = projectsApiImpl;
export const epicsApi = epicsApiImpl;
export const boardsApi = boardsApiImpl;
export const tasksApi = tasksApiImpl;
export const labelsApi = labelsApiImpl;
export const wikiApi = wikiApiImpl;
export const filesApi = filesApiImpl;
export const adminApi = adminApiImpl;

export * from './types';
