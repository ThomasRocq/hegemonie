{% include "header.tpl" %}
<div class="large"><h2>Stock</h2>
    <table class="budget">
        <thead>
        <tr>
            <td class="title">Now</td>
            <td>{{Land.Stock.Usage.R0}}</td>
            <td>{{Land.Stock.Usage.R1}}</td>
            <td>{{Land.Stock.Usage.R2}}</td>
            <td>{{Land.Stock.Usage.R3}}</td>
            <td>{{Land.Stock.Usage.R4}}</td>
            <td>{{Land.Stock.Usage.R5}}</td>
        </tr>
        </thead>
        <tbody>
        <tr>
            <td class="title">Base Capacity</td>
            <td>{{Land.Stock.Base.R0}}</td>
            <td>{{Land.Stock.Base.R1}}</td>
            <td>{{Land.Stock.Base.R2}}</td>
            <td>{{Land.Stock.Base.R3}}</td>
            <td>{{Land.Stock.Base.R4}}</td>
            <td>{{Land.Stock.Base.R5}}</td>
        </tr>
        <tr>
            <td class="title">Buildings</td>
            <td>{{Land.Stock.Buildings.Mult.R0}}</td>
            <td>{{Land.Stock.Buildings.Mult.R1}}</td>
            <td>{{Land.Stock.Buildings.Mult.R2}}</td>
            <td>{{Land.Stock.Buildings.Mult.R3}}</td>
            <td>{{Land.Stock.Buildings.Mult.R4}}</td>
            <td>{{Land.Stock.Buildings.Mult.R5}}</td>
        </tr>
        <tr>
            <td class="title">Knowledge</td>
            <td>{{Land.Stock.Knowledge.Mult.R0}}</td>
            <td>{{Land.Stock.Knowledge.Mult.R1}}</td>
            <td>{{Land.Stock.Knowledge.Mult.R2}}</td>
            <td>{{Land.Stock.Knowledge.Mult.R3}}</td>
            <td>{{Land.Stock.Knowledge.Mult.R4}}</td>
            <td>{{Land.Stock.Knowledge.Mult.R5}}</td>
        </tr>
        <tr>
            <td class="title">Troops</td>
            <td>{{Land.Stock.Troops.Mult.R0}}</td>
            <td>{{Land.Stock.Troops.Mult.R1}}</td>
            <td>{{Land.Stock.Troops.Mult.R2}}</td>
            <td>{{Land.Stock.Troops.Mult.R3}}</td>
            <td>{{Land.Stock.Troops.Mult.R4}}</td>
            <td>{{Land.Stock.Troops.Mult.R5}}</td>
        </tr>
        <tr>
            <td class="title">Buildings</td>
            <td>{{Land.Stock.Buildings.Plus.R0}}</td>
            <td>{{Land.Stock.Buildings.Plus.R1}}</td>
            <td>{{Land.Stock.Buildings.Plus.R2}}</td>
            <td>{{Land.Stock.Buildings.Plus.R3}}</td>
            <td>{{Land.Stock.Buildings.Plus.R4}}</td>
            <td>{{Land.Stock.Buildings.Plus.R5}}</td>
        </tr>
        <tr>
            <td class="title">Knowledge</td>
            <td>{{Land.Stock.Knowledge.Plus.R0}}</td>
            <td>{{Land.Stock.Knowledge.Plus.R1}}</td>
            <td>{{Land.Stock.Knowledge.Plus.R2}}</td>
            <td>{{Land.Stock.Knowledge.Plus.R3}}</td>
            <td>{{Land.Stock.Knowledge.Plus.R4}}</td>
            <td>{{Land.Stock.Knowledge.Plus.R5}}</td>
        </tr>
        <tr>
            <td class="title">Troops</td>
            <td>{{Land.Stock.Troops.Plus.R0}}</td>
            <td>{{Land.Stock.Troops.Plus.R1}}</td>
            <td>{{Land.Stock.Troops.Plus.R2}}</td>
            <td>{{Land.Stock.Troops.Plus.R3}}</td>
            <td>{{Land.Stock.Troops.Plus.R4}}</td>
            <td>{{Land.Stock.Troops.Plus.R5}}</td>
        </tr>
        </tbody>
        <tfoot>
        <tr>
            <td class="title">Max</td>
            <td>{{Land.Stock.Actual.R0}}</td>
            <td>{{Land.Stock.Actual.R1}}</td>
            <td>{{Land.Stock.Actual.R2}}</td>
            <td>{{Land.Stock.Actual.R3}}</td>
            <td>{{Land.Stock.Actual.R4}}</td>
            <td>{{Land.Stock.Actual.R5}}</td>
        </tr>
        </tfoot>
    </table>
</div>
<div class="large"><h2>Production</h2>
    <table class="budget">
        <thead>
        <tr>
            <td class="title">Total</td>
            <td>{{Land.Production.Actual.R0}}</td>
            <td>{{Land.Production.Actual.R1}}</td>
            <td>{{Land.Production.Actual.R2}}</td>
            <td>{{Land.Production.Actual.R3}}</td>
            <td>{{Land.Production.Actual.R4}}</td>
            <td>{{Land.Production.Actual.R5}}</td>
        </tr>
        </thead>
        <tbody>
        <tr>
            <td class="title">Base Production</td>
            <td>{{Land.Production.Base.R0}}</td>
            <td>{{Land.Production.Base.R1}}</td>
            <td>{{Land.Production.Base.R2}}</td>
            <td>{{Land.Production.Base.R3}}</td>
            <td>{{Land.Production.Base.R4}}</td>
            <td>{{Land.Production.Base.R5}}</td>
        </tr>
        <tr>
            <td class="title">Buildings</td>
            <td>{{Land.Production.Buildings.Mult.R0}}</td>
            <td>{{Land.Production.Buildings.Mult.R1}}</td>
            <td>{{Land.Production.Buildings.Mult.R2}}</td>
            <td>{{Land.Production.Buildings.Mult.R3}}</td>
            <td>{{Land.Production.Buildings.Mult.R4}}</td>
            <td>{{Land.Production.Buildings.Mult.R5}}</td>
        </tr>
        <tr>
            <td class="title">Knowledge</td>
            <td>{{Land.Production.Knowledge.Mult.R0}}</td>
            <td>{{Land.Production.Knowledge.Mult.R1}}</td>
            <td>{{Land.Production.Knowledge.Mult.R2}}</td>
            <td>{{Land.Production.Knowledge.Mult.R3}}</td>
            <td>{{Land.Production.Knowledge.Mult.R4}}</td>
            <td>{{Land.Production.Knowledge.Mult.R5}}</td>
        </tr>
        <tr>
            <td class="title">Troops</td>
            <td>{{Land.Production.Troops.Mult.R0}}</td>
            <td>{{Land.Production.Troops.Mult.R1}}</td>
            <td>{{Land.Production.Troops.Mult.R2}}</td>
            <td>{{Land.Production.Troops.Mult.R3}}</td>
            <td>{{Land.Production.Troops.Mult.R4}}</td>
            <td>{{Land.Production.Troops.Mult.R5}}</td>
        </tr>
        <tr>
            <td class="title">Buildings</td>
            <td>{{Land.Production.Buildings.Plus.R0}}</td>
            <td>{{Land.Production.Buildings.Plus.R1}}</td>
            <td>{{Land.Production.Buildings.Plus.R2}}</td>
            <td>{{Land.Production.Buildings.Plus.R3}}</td>
            <td>{{Land.Production.Buildings.Plus.R4}}</td>
            <td>{{Land.Production.Buildings.Plus.R5}}</td>
        </tr>
        <tr>
            <td class="title">Knowledge</td>
            <td>{{Land.Production.Knowledge.Plus.R0}}</td>
            <td>{{Land.Production.Knowledge.Plus.R1}}</td>
            <td>{{Land.Production.Knowledge.Plus.R2}}</td>
            <td>{{Land.Production.Knowledge.Plus.R3}}</td>
            <td>{{Land.Production.Knowledge.Plus.R4}}</td>
            <td>{{Land.Production.Knowledge.Plus.R5}}</td>
        </tr>
        <tr>
            <td class="title">Troops</td>
            <td>{{Land.Production.Troops.Plus.R0}}</td>
            <td>{{Land.Production.Troops.Plus.R1}}</td>
            <td>{{Land.Production.Troops.Plus.R2}}</td>
            <td>{{Land.Production.Troops.Plus.R3}}</td>
            <td>{{Land.Production.Troops.Plus.R4}}</td>
            <td>{{Land.Production.Troops.Plus.R5}}</td>
        </tr>
        </tbody>
        <tfoot>
        </tfoot>
    </table>
</div>
{% include "footer.tpl" %}
