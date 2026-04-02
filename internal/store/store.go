package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Goal struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Area string `json:"area"`
	Target string `json:"target"`
	Progress int `json:"progress"`
	DueDate string `json:"due_date"`
	Status string `json:"status"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"meridian.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS goals(id TEXT PRIMARY KEY,title TEXT NOT NULL,area TEXT DEFAULT '',target TEXT DEFAULT '',progress INTEGER DEFAULT 0,due_date TEXT DEFAULT '',status TEXT DEFAULT 'active',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Goal)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO goals(id,title,area,target,progress,due_date,status,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Area,e.Target,e.Progress,e.DueDate,e.Status,e.Notes,e.CreatedAt);return err}
func(d *DB)Get(id string)*Goal{var e Goal;if d.db.QueryRow(`SELECT id,title,area,target,progress,due_date,status,notes,created_at FROM goals WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Area,&e.Target,&e.Progress,&e.DueDate,&e.Status,&e.Notes,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Goal{rows,_:=d.db.Query(`SELECT id,title,area,target,progress,due_date,status,notes,created_at FROM goals ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Goal;for rows.Next(){var e Goal;rows.Scan(&e.ID,&e.Title,&e.Area,&e.Target,&e.Progress,&e.DueDate,&e.Status,&e.Notes,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Goal)error{_,err:=d.db.Exec(`UPDATE goals SET title=?,area=?,target=?,progress=?,due_date=?,status=?,notes=? WHERE id=?`,e.Title,e.Area,e.Target,e.Progress,e.DueDate,e.Status,e.Notes,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM goals WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM goals`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Goal{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (title LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,title,area,target,progress,due_date,status,notes,created_at FROM goals WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Goal;for rows.Next(){var e Goal;rows.Scan(&e.ID,&e.Title,&e.Area,&e.Target,&e.Progress,&e.DueDate,&e.Status,&e.Notes,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM goals GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
