import { BaseAPI, route } from "../base";
import type {
  SetupCreate,
  SetupOut,
  SetupPartAddOut,
  SetupPartCreate,
  SetupPartOut,
  SetupStatsOut,
  SetupUpdate,
  UsageSessionCreate,
  UsageSessionOut,
  UsageSessionUpdate,
} from "../types/data-contracts";

export class SetupsAPI extends BaseAPI {
  getAll() {
    return this.http.get<SetupOut[]>({ url: route(`/setups`) });
  }

  create(data: SetupCreate) {
    return this.http.post<SetupCreate, SetupOut>({ url: route(`/setups`), body: data });
  }

  get(id: string) {
    return this.http.get<SetupOut>({ url: route(`/setups/${id}`) });
  }

  update(id: string, data: SetupUpdate) {
    return this.http.put<SetupUpdate, SetupOut>({ url: route(`/setups/${id}`), body: data });
  }

  delete(id: string) {
    return this.http.delete<void>({ url: route(`/setups/${id}`) });
  }

  getParts(id: string) {
    return this.http.get<SetupPartOut[]>({ url: route(`/setups/${id}/parts`) });
  }

  getCurrentParts(id: string) {
    return this.http.get<SetupPartOut[]>({ url: route(`/setups/${id}/parts/current`) });
  }

  addPart(setupId: string, data: SetupPartCreate) {
    return this.http.post<SetupPartCreate, SetupPartAddOut>({
      url: route(`/setups/${setupId}/parts`),
      body: data,
    });
  }

  getStats() {
    return this.http.get<SetupStatsOut>({ url: route(`/setups/statistics`) });
  }

  removePart(setupId: string, partId: string) {
    return this.http.delete<void>({ url: route(`/setups/${setupId}/parts/${partId}`) });
  }

  purgePart(setupId: string, partId: string) {
    return this.http.delete<void>({ url: route(`/setups/${setupId}/parts/${partId}/purge`) });
  }

  getSessions(setupId: string) {
    return this.http.get<UsageSessionOut[]>({ url: route(`/setups/${setupId}/sessions`) });
  }

  createSession(setupId: string, data: UsageSessionCreate) {
    return this.http.post<UsageSessionCreate, UsageSessionOut>({
      url: route(`/setups/${setupId}/sessions`),
      body: data,
    });
  }

  updateSession(setupId: string, sessionId: string, data: UsageSessionUpdate) {
    return this.http.put<UsageSessionUpdate, UsageSessionOut>({
      url: route(`/setups/${setupId}/sessions/${sessionId}`),
      body: data,
    });
  }

  deleteSession(setupId: string, sessionId: string) {
    return this.http.delete<void>({ url: route(`/setups/${setupId}/sessions/${sessionId}`) });
  }

  getEntityHistory(entityId: string) {
    return this.http.get<SetupPartOut[]>({ url: route(`/entities/${entityId}/setup-history`) });
  }
}
