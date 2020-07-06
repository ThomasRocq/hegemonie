{% include "header.tpl" %}
<script src="/static/hege-map.js"></script>
<script>
// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

const onField = function (field, v) { if (field != null) { field.value = v; } };

const onForm = function (form, pos, cityId) {
  let items = form.children;
  for (let i=0; i<items.length; i++) {
    let f = items[i];
    if (f.id === "Location") {
      f.value = pos;
    } else if (f.id === "CityID") {
      f.value = cityId;
    }
  }
};

const allForms = function (doc, pos, cityId, cityName) {
  doc.getElementById("CityName").value = cityName;
  let forms = doc.getElementsByTagName('form');
  for (let i = 0; i < forms.length; i++) {
    onForm(forms[i], pos, cityId);
  }
};

window.addEventListener("load", function() {
  let svg1 = document.getElementById('interactive-map');
  let armies = [{"id":{{aid}}, "cell":{{Army.Location}}}];
  drawMapWithCities(svg1, "calaquyr",
    function (position) {
      allForms(document, pos, null, null);
    },
    function (position, cityId, cityName) {
      allForms(document, position, cityId, cityName);
    })
    .then(map => {
      hightlightCell(svg1, {{Land.Location}});
      return map;
    })
    .then(map => {
      return patchWithArmies(svg1, map, armies);
    })
    .catch(err => { console.log(err); });
});
</script>
{% include "map.tpl" %}

<div><h2>Actions</h2>

<p>Target: <input type="text" id="CityName" value=""/></p>

<form if="action-move" method="post" action="/action/army/move">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="location" id="Location" value=""/>
    <input type="submit" value="Move There"/>
</form>

<form if="action-wait" method="post" action="/action/army/wait">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="location" id="Location" value=""/>
    <input type="submit" value="Wait There"/>
</form>

<form if="action-defend" method="post" action="/action/army/defend">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="location" id="Location" value=""/>
    <input type="submit" value="Defend"/>
</form>

<form if="action-assault" method="post" action="/action/army/assault">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="location" id="Location" value=""/>
    <input type="submit" value="Assault"/>
</form>

<form if="action-disband" method="post" action="/action/army/cancel">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="submit" value="Disband There"/>
</form>
<br/>
<form if="action-disband" method="post" action="/action/army/flea">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="submit" value="Flea from the fight"/>
</form>

<form if="action-disband" method="post" action="/action/army/flip">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="submit" value="Flip in the fight"/>
</form>

<form if="action-disband" method="post" action="/action/army/cancel">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="submit" value="Dismantle Here"/>
</form>

</div>

<div><h2>Commands</h2>
{% for cmd in Commands %}
    <p>{{cmd.SeqNum}}:
    {% if cmd.CommandID == 0 %}
    Take selfies on
    {% elif cmd.CommandID == 1 %}
    Disband on
    {% elif cmd.CommandID == 2 %}
    Wait on
    {% elif cmd.CommandID == 3 %}
    Move to
    {% elif cmd.CommandID == 4 %}
    Assault
    {% elif cmd.CommandID == 5 %}
    Go defend
    {% else %}
    ?
    {% endif %} (type {{cmd.CommandID}})
    {{cmd.CityName}} (id {{cmd.CityID}}, located at {{cmd.Location}})</p>
    <form method="post" action="/action/army/command/top">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="command" value="{{cmd.SeqNum}}"/>
    <input type="submit" value="Top"/>
    </form>
    <form method="post" action="/action/army/command/drop">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="command" value="{{cmd.SeqNum}}"/>
    <input type="submit" value="Drop"/>
    </form>
{% endfor %}

    <form method="post" action="/action/army/command/flush">
    <input type="hidden" name="cid" value="{{Character.Id}}"/>
    <input type="hidden" name="lid" value="{{Land.Id}}"/>
    <input type="hidden" name="aid" value="{{Army.Id}}"/>
    <input type="hidden" name="command" value="{{cmd.SeqNum}}"/>
    <input type="submit" value="Top"/>
    </form>

</div>

<div><h2>Payload</h2>
    <p>
    Gold ({{Army.Stock.R1}}),
    Cereals ({{Army.Stock.R2}}),
    Livestock ({{Army.Stock.R3}}),
    Wood ({{Army.Stock.R4}}),
    Stone ({{Army.Stock.R5}})
    </p>
</div>

<div><h2>Enroll</h2>
    {% for u in Land.Assets.Units %}
    {% if u.Ticks == 0 %}
    <p>{{u.Type.Name}} (id {{u.Id}}) Health({{u.Health}}/{{u.Type.Health}})</p>
    {% endif %}
    {% endfor %}
</div>

<div><h2>Units</h2>
    {% for u in Army.Units %}
    <p>{{u.Type.Name}} (id {{u.Id}}) Health({{u.Health}}/{{u.Type.Health}})</p>
    {% endfor %}
</div>

{% include "footer.tpl" %}
