library Shine initializer Init requires NLBonus
globals
	public boolexpr Ally
	public boolexpr Enemy
	public timer Cycle = CreateTimer()
	public group All = CreateGroup()
	public integer array Num
	public real array PosX
	public real array PosY
endglobals

public function ally takes nothing returns boolean
	return IsUnitAlly(GetFilterUnit(), Player(0))
endfunction
public function enemy takes nothing returns boolean
	return IsUnitEnemy(GetFilterUnit(), Player(0))
endfunction


public function Register takes unit u,integer buf,real range, boolexpr filter returns integer
	call TriggerRegisterUnitInRange(tri,u,range,filter)
endfunction

private function setpos takes nothing returns nothing
	set PosX[GetUnitUserData(GetEnumUnit())] = GetUnitX(GetEnumUnit())
	set PosX[GetUnitUserData(GetEnumUnit())] = GetUnitX(GetEnumUnit())
endfunction
private function setpos takes nothing returns nothing
	ForGroup(All,function setpos)
endfunction

private function Init takes nothing returns nothing
	call TimerStart(Cycle,.8,true,function setpos)
	set Ally = boolexpr(function ally)
	set Enemy = boolexpr(function enemy)
endfunction

endlibrary