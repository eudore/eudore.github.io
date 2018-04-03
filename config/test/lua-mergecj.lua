
function table.clone(object)
    local lookup_table = {}
    local function _copy(object)
        if type(object) ~= "table" then
            return object
        elseif lookup_table[object] then
            return lookup_table[object]
        end
        local new_table = {}
        lookup_table[object] = new_table
        for index, value in pairs(object) do
            new_table[_copy(index)] = _copy(value)
        end
        return setmetatable(new_table, getmetatable(object))
    end
    return _copy(object)
end


key = ""
function PrintTable(table , level)
	level = level or 1
	local indent = ""
	for i = 1, level do
		indent = indent.."  "
	end

	if key ~= "" then
		print(indent..key.." ".."=".." ".."{")
	else
		print(indent .. "{")
	end

	key = ""
	for k,v in pairs(table) do
		if type(v) == "table" then
			key = k
			PrintTable(v, level + 1)
		else
			local content = string.format("%s%s = %s", indent .. "  ",tostring(k), tostring(v))
			print(content)  
		end
	end
	print(indent .. "}")

end

function table.len(tab)
	if tab == nil then
		return 0
	end
	local i = 0
	for _,_ in pairs(tab) do
		i=i+1
	end
	return i
end


function JsMatch(data)
	if table.len(data) == 0 then
		return ""
	elseif table.len(data) == 1 then
		for k,v in pairs(data) do
			if type(v) ~= "table" then
				return k
			end
		end
	end
	local d = "(?:"
	for k,v in pairs(data) do
		if type(v) == "table" then
			d = d..k..JsMatch(v).."|"
		else
			d = d..k .."|"
		end
	end
	return string.gsub(string.gsub(d,"@",""),"|$",")")
	--return string.gsub(d,"|$",")")
end

function Branch(data)
	local new={}
	for k,v in pairs(data) do
		if type(v) ~= "table" then
			local first = string.sub(v,1,1)
			local other = string.sub(v,2,-1)
			if new[first] == nil then
				new[first]={}
			end
			if string.len(v) > 1 then
				new[first][other] = other
			end
		end
	end
	for k,v in pairs(new) do
		if type(v) == "table" then
			local result = Branch(v)
			if result ~= nil then
				new[k]=result
			end
		end
	end
	local n = table.len(new)
	if n == 0 then
		return nil
	else
		return new
	end
end

function Merge(data)
	local new = {}
	for k,v in pairs(data) do
		local n,_ =string.find(k,"@")
		if n~=nil then
			new[k] ={["@"]="@"}
		elseif type(v) == "table" then
			local r = Merge(v)
			if table.len(r) == 1 then
				for vk,vv in pairs(r) do
					new[k..vk]=vv
				end
			else
				new[k]=r
			end
		end
	end
	return new
end


local data={}
local file=io.open("data-cj.txt","r");
for line in file:lines() do
	data[string.gsub(line.."@"," ","")]=string.gsub(line.."@"," ","")
end


str=JsMatch(data)
print(string.len(str))


--PrintTable(data)
data=Branch(data)
--PrintTable(data)
data=Merge(data)
data=Merge(data)
--PrintTable(data)
str=JsMatch(data)
print(str)
print(string.len(str))

--12258/19038