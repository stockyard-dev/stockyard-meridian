package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-meridian/internal/store")
func(s *Server)handleListWidgets(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListWidgets();if list==nil{list=[]store.Widget{}};writeJSON(w,200,list)}
func(s *Server)handleCreateWidget(w http.ResponseWriter,r *http.Request){var wid store.Widget;json.NewDecoder(r.Body).Decode(&wid);if wid.Name==""{writeError(w,400,"name required");return};wid.Enabled=true;if wid.Config==""{wid.Config="{}"};s.db.CreateWidget(&wid);writeJSON(w,201,wid)}
func(s *Server)handleUpdateWidget(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Value string `json:"value"`};json.NewDecoder(r.Body).Decode(&req);s.db.UpdateWidget(id,req.Value);writeJSON(w,200,map[string]string{"status":"updated"})}
func(s *Server)handleLogEntry(w http.ResponseWriter,r *http.Request){var e store.LifeEntry;json.NewDecoder(r.Body).Decode(&e);if e.Category==""||e.Title==""{writeError(w,400,"category and title required");return};s.db.LogEntry(&e);writeJSON(w,201,e)}
func(s *Server)handleListEntries(w http.ResponseWriter,r *http.Request){cat:=r.URL.Query().Get("category");list,_:=s.db.ListEntries(cat);if list==nil{list=[]store.LifeEntry{}};writeJSON(w,200,list)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
