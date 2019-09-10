package slog

// region --- INFO Level Sugars ---
// Note logs out a message in INFO level and with Operation NOTE. Returns an instance of operation NOTE
func (i *slogInstance) Note(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(NOTE).Info(str, v...)
}

// Await logs out a message in INFO level and with Operation AWAIT. Returns an instance of operation AWAIT
func (i *slogInstance) Await(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(AWAIT).Info(str, v...)
}

// Done logs out a message in INFO level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) Done(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Info(str, v...)
}

// Success logs out a message in INFO level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) Success(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Info(str, v...)
}

// IO logs out a message in INFO level and with Operation IO. Returns an instance of operation IO
func (i *slogInstance) IO(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(IO).Info(str, v...)
}

// endregion
// region --- WARN Level Sugars ---
// Note logs out a message in WARN level and with Operation NOTE. Returns an instance of operation NOTE
func (i *slogInstance) WarnNote(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(NOTE).Warn(str, v...)
}

// Await logs out a message in WARN level and with Operation AWAIT. Returns an instance of operation AWAIT
func (i *slogInstance) WarnAwait(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(AWAIT).Warn(str, v...)
}

// Done logs out a message in WARN level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) WarnDone(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Warn(str, v...)
}

// Success logs out a message in WARN level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) WarnSuccess(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Warn(str, v...)
}

// IO logs out a message in WARN level and with Operation IO. Returns an instance of operation IO
func (i *slogInstance) WarnIO(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(IO).Warn(str, v...)
}

// endregion
// region --- ERROR Level Sugars ---
// Note logs out a message in ERROR level and with Operation NOTE. Returns an instance of operation NOTE
func (i *slogInstance) ErrorNote(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(NOTE).Error(str, v...)
}

// Await logs out a message in ERROR level and with Operation AWAIT. Returns an instance of operation AWAIT
func (i *slogInstance) ErrorAwait(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(AWAIT).Error(str, v...)
}

// Done logs out a message in ERROR level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) ErrorDone(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Error(str, v...)
}

// Success logs out a message in ERROR level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) ErrorSuccess(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Error(str, v...)
}

// IO logs out a message in ERROR level and with Operation IO. Returns an instance of operation IO
func (i *slogInstance) ErrorIO(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(IO).Error(str, v...)
}

// endregion
// region --- DEBUG Level Sugars ---
// Note logs out a message in DEBUG level and with Operation NOTE. Returns an instance of operation NOTE
func (i *slogInstance) DebugNote(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(NOTE).Debug(str, v...)
}

// Await logs out a message in DEBUG level and with Operation AWAIT. Returns an instance of operation AWAIT
func (i *slogInstance) DebugAwait(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(AWAIT).Debug(str, v...)
}

// Done logs out a message in DEBUG level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) DebugDone(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Debug(str, v...)
}

// Success logs out a message in DEBUG level and with Operation DONE. Returns an instance of operation DONE
func (i *slogInstance) DebugSuccess(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(DONE).Debug(str, v...)
}

// IO logs out a message in DEBUG level and with Operation IO. Returns an instance of operation IO
func (i *slogInstance) DebugIO(str interface{}, v ...interface{}) Instance {
	return i.clone().incStackOffset().Operation(IO).Debug(str, v...)
}

// endregion
