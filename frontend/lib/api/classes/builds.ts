import { BaseAPI, route } from "../base";
import type {
  BuildCreate,
  BuildOut,
  BuildPartAddOut,
  BuildPartCreate,
  BuildPartOut,
  BuildStatsOut,
  SportPositionConfig,
  BuildUpdate,
  UsageSessionCreate,
  UsageSessionOut,
  UsageSessionUpdate,
} from "../types/data-contracts";

export class BuildsAPI extends BaseAPI {
  getSports() {
    return this.http.get<SportPositionConfig[]>({ url: route(`/builds/sports`) });
  }

  getAll() {
    return this.http.get<BuildOut[]>({ url: route(`/builds`) });
  }

  create(data: BuildCreate) {
    return this.http.post<BuildCreate, BuildOut>({ url: route(`/builds`), body: data });
  }

  get(id: string) {
    return this.http.get<BuildOut>({ url: route(`/builds/${id}`) });
  }

  update(id: string, data: BuildUpdate) {
    return this.http.put<BuildUpdate, BuildOut>({ url: route(`/builds/${id}`), body: data });
  }

  delete(id: string) {
    return this.http.delete<void>({ url: route(`/builds/${id}`) });
  }

  getParts(id: string) {
    return this.http.get<BuildPartOut[]>({ url: route(`/builds/${id}/parts`) });
  }

  getCurrentParts(id: string) {
    return this.http.get<BuildPartOut[]>({ url: route(`/builds/${id}/parts/current`) });
  }

  addPart(buildId: string, data: BuildPartCreate) {
    return this.http.post<BuildPartCreate, BuildPartAddOut>({
      url: route(`/builds/${buildId}/parts`),
      body: data,
    });
  }

  getStats() {
    return this.http.get<BuildStatsOut>({ url: route(`/builds/statistics`) });
  }

  removePart(buildId: string, partId: string) {
    return this.http.delete<void>({ url: route(`/builds/${buildId}/parts/${partId}`) });
  }

  purgePart(buildId: string, partId: string) {
    return this.http.delete<void>({ url: route(`/builds/${buildId}/parts/${partId}/purge`) });
  }

  getSessions(buildId: string) {
    return this.http.get<UsageSessionOut[]>({ url: route(`/builds/${buildId}/sessions`) });
  }

  createSession(buildId: string, data: UsageSessionCreate) {
    return this.http.post<UsageSessionCreate, UsageSessionOut>({
      url: route(`/builds/${buildId}/sessions`),
      body: data,
    });
  }

  updateSession(buildId: string, sessionId: string, data: UsageSessionUpdate) {
    return this.http.put<UsageSessionUpdate, UsageSessionOut>({
      url: route(`/builds/${buildId}/sessions/${sessionId}`),
      body: data,
    });
  }

  deleteSession(buildId: string, sessionId: string) {
    return this.http.delete<void>({ url: route(`/builds/${buildId}/sessions/${sessionId}`) });
  }

  getEntityHistory(entityId: string) {
    return this.http.get<BuildPartOut[]>({ url: route(`/entities/${entityId}/build-history`) });
  }
}
