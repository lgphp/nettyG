package netgo

import "io/ioutil"

type HandlerContext struct {
	p *Pipeline
    next *HandlerContext
	prev *HandlerContext
	handler Handler
}

func newHandlerContext(p *Pipeline,handler Handler) *HandlerContext {
	return &HandlerContext{p:p,handler:handler}
}

func (h *HandlerContext) FireWrite(data interface{}){
	hc:=h.findNextOutbound()
	if hc!=nil{
		hc.handler.(OutboundHandler).Write(hc,data)
	}
}

func (h *HandlerContext) FireRead(data interface{}){
    hc:=h.findNextInbound()
	if hc!=nil{
		hc.handler.(InboundHandler).Read(hc,data)
	}
}

func (h *HandlerContext) FireConnected(){
	hc:=h.findNextInbound()
	if hc!=nil{
		hc.handler.(InboundHandler).Connected(hc)
	}
}

func (h *HandlerContext) isInbound() bool{
	_,ok:= h.handler.(InboundHandler)
	return  ok
}

func (h *HandlerContext) isOutbound() bool{
	_,ok:= h.handler.(OutboundHandler)
	return ok
}

func (h *HandlerContext) findNextInbound() *HandlerContext{
	next:=h
	for{
		next=next.next

		if next.isInbound(){
			return next
		}
	}
}

func (h *HandlerContext) findNextOutbound() *HandlerContext{
	prev:=h
	for{

		prev=prev.prev

		if prev.isOutbound(){
			return prev
		}
	}
}


func (h *HandlerContext) WriteToReadBuffer(b []byte){
	h.p.chl.readbuffer.Write(b)
}

func (h *HandlerContext) ResetReadBuffer(){
	h.p.chl.readbuffer.Reset()
}

func (h *HandlerContext) ReadAllReadBuffer() ([]byte,error){
	return ioutil.ReadAll(h.p.chl.readbuffer)
}

