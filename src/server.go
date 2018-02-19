package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

var hubs map[string]*Hub

func init() {
	os.MkdirAll("data", 0755)
}

// Run will run the main program
func Run(port string) (err error) {
	defer log.Flush()

	hubs = make(map[string]*Hub)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			namesToDelete := make(map[string]struct{})
			for name := range hubs {
				// log.Debugf("hub %s has %d clients", name, len(hubs[name].clients))
				if len(hubs[name].clients) == 0 {
					namesToDelete[name] = struct{}{}
				}
			}
			for name := range namesToDelete {
				log.Debugf("deleting hub for %s", name)
				delete(hubs, name)
			}
		}
	}()

	// load static stuff
	mainCSS, err := ioutil.ReadFile(path.Join("static", "tachyons.min.css"))
	if err != nil {
		return
	}
	// setup gin server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	// Standardize logs
	r.LoadHTMLGlob("templates/*")
	r.Use(middleWareHandler(), gin.Recovery())
	r.HEAD("/", func(c *gin.Context) { // handler for the uptime robot
		c.String(http.StatusOK, "OK")
	})
	r.GET("/*name", func(c *gin.Context) {
		name := c.Param("name")
		if len(name) == 1 {
			c.String(http.StatusOK, "OK")
		} else if name == "/ws" {
			name = c.DefaultQuery("name", "")
			if name == "" {
				c.String(http.StatusOK, "OK")
				return
			}
			if _, ok := hubs[name]; !ok {
				hubs[name] = newHub(name)
				go hubs[name].run()
				time.Sleep(50 * time.Millisecond)
			}
			hubs[name].serveWs(c.Writer, c.Request)
		} else if strings.Contains(name, "/static") {
			c.Data(http.StatusOK, "text/css", mainCSS)
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Name": name[1:],
			})
		}
	})
	r.POST("/", handlerPostMessage) // handle for posting message
	log.Infof("Running at http://0.0.0.0:" + port)
	err = r.Run(":" + port)
	return
}

func handlerPostMessage(c *gin.Context) {
	message, err := func(c *gin.Context) (message string, err error) {
		var m messageJSON
		err = c.ShouldBindJSON(&m)
		if err != nil {
			return
		}
		log.Debug("bound message")
		message = fmt.Sprintf("got message for %s", m.To)
		m, err = validateMessage(m)
		if err != nil {
			return
		}
		log.Debug("validated messages")
		db := open(m.To)
		err = db.saveMessage(m)
		if err != nil {
			log.Error(err)
		}
		db.close()
		log.Debug("saved message")
		if _, ok := hubs[m.To]; ok {
			hubs[m.To].broadcastNextMessage(false)
		}
		log.Debug("broadcast message")
		return
	}(c)
	if err != nil {
		message = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"success": err == nil,
	})
}

func middleWareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// Add base headers
		addCORS(c)
		// Run next function
		c.Next()
		// Log request
		log.Infof("%v %v %v %s", c.Request.RemoteAddr, c.Request.Method, c.Request.URL, time.Since(t))
	}
}

func addCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
}
